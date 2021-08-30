// Copyright 2021 tobyzxj
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package netproxy

import (
	"errors"
	"io"
	"net"
	"time"
)

// TCPHandlerCallback 用户定义的数据处理回调函数
type TCPHandlerCallback func(conn *net.TCPConn, addr net.Addr, data []byte, length int) error

// TCP 服务器
type TCProxy struct {
	Name          string             // 服务器名称定义
	Addr          string             // TCP服务器监听地址
	Listener      *net.TCPListener   // TCP 服务器
	Conn          *net.TCPConn       // TCP连接
	ReadDeadline  int                // 读取超时, 单位秒
	WriteDeadline int                // 发送超时，单位秒
	Handler       TCPHandlerCallback // 用户自定义处理回调函数
}

// Start 启动一个TCP服务器
func (u *TCProxy) Start() error {
	if u != nil {
		if u.Addr == "" {
			return errors.New("addr is empty")
		}
		tcpAddr, err := net.ResolveTCPAddr("tcp4", u.Addr)
		if err != nil {
			return err
		}
		listener, err := net.ListenTCP("tcp4", tcpAddr)
		if err != nil {
			TraceError("[tcproxy] listen tcp port failed, %v", err)
			return err
		}
		defer listener.Close()
		u.Listener = listener
		if u.Handler != nil {
			handleTCPPacket(u.Listener, u.Handler, u.ReadDeadline, u.WriteDeadline)
		} else {
			TraceError("[tcproxy] define user-specified callback function first")
			return errors.New("define user-specified callback function first")
		}
		// never return
		return nil
	}
	return errors.New("u is nil")
}

// handleTCPPacket Handle tcp net packet
func handleTCPPacket(listener *net.TCPListener, handle TCPHandlerCallback, readDeadline, writeDeadline int) {
	// loop forever
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			TraceError("[tcproxy] tcp accept failed, %v", err)
			continue
		}

		// delay response data
		buf := make([]byte, 4096)
		go handleTCPPacketThread(conn, conn.RemoteAddr(), buf, handle, readDeadline, writeDeadline)
	}
}

// handleTCPPacketThread Response client's packet
func handleTCPPacketThread(conn *net.TCPConn, addr net.Addr, data []byte, handle TCPHandlerCallback, readDeadline, writeDeadline int) error {
	defer func() {
		if conn != nil {
			TraceInfo("[tcproxy] tcp client close: <- %s", addr)
			conn.Close()
		}

		// recover panic
		if err := recover(); err != nil {
			TraceError("[tcproxy] handle packet panic, %v", err)
		}

		data = nil
	}()

	TraceInfo("[tcproxy] new tcp connected: -> %s", addr)
	for {
		// set read timeout
		tRead := time.Now().Add(time.Minute * time.Duration(readDeadline))
		tWrite := time.Now().Add(time.Minute * time.Duration(writeDeadline))
		conn.SetReadDeadline(tRead)
		conn.SetWriteDeadline(tWrite)

		// start to read
		length, err := conn.Read(data[:])
		if err != nil {
			TraceError("[tcproxy] tcp read data failed, %v", err)
			if err == io.EOF {
				return err
			}
			return err
		}
		if handle != nil {
			err = handle(conn, addr, data, length)
			if err != nil {
				return err
			}
		}
	}
}
