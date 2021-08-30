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
	"net"
)

// UDPHandlerCallback 用户定义的数据处理回调函数
type UDPHandlerCallback func(conn *net.UDPConn, addr *net.UDPAddr, data []byte, length int) error

// UDP 服务器
type UDProxy struct {
	Name    string             // 服务器名称定义
	Addr    string             // UDP服务器监听地址
	Conn    *net.UDPConn       // UDP连接
	Handler UDPHandlerCallback // 用户自定义处理回调函数
}

// Start 启动一个UDP服务器
func (u *UDProxy) Start() error {
	if u != nil {
		if u.Addr == "" {
			return errors.New("addr is empty")
		}
		udpAddr, err := net.ResolveUDPAddr("udp4", u.Addr)
		if err != nil {
			return err
		}
		conn, err := net.ListenUDP("udp4", udpAddr)
		if err != nil {
			TraceError("[udproxy] listen udp port failed, %v", err)
			return err
		}
		defer conn.Close()
		u.Conn = conn
		if u.Handler != nil {
			handleUDPPacket(u.Conn, u.Handler)
		} else {
			TraceError("[udproxy] define user-specified callback function first")
			return errors.New("define user-specified callback function first")
		}
		// never return
		return nil
	}
	return errors.New("u is nil")
}

// handleUDPPacket Handle udp net packet
func handleUDPPacket(conn *net.UDPConn, handle UDPHandlerCallback) {
	// loop forever
	for {
		// read data
		buf := make([]byte, 4096)
		// ReadFromUDP 阻塞的
		num, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			TraceError("[udproxy] udp receive failed, %v", err)
			continue
		}

		// delay response data
		go handleUDPPacketThread(conn, addr, buf, num, handle)
	}

}

// handleUDPPacketThread Response client's packet
func handleUDPPacketThread(conn *net.UDPConn, addr *net.UDPAddr, data []byte, length int, handle UDPHandlerCallback) error {
	var err error

	defer func() {
		data = nil
		// recover panic
		if err := recover(); err != nil {
			TraceError("[udproxy] handle packet panic, %v", err)
		}
	}()

	if handle != nil {
		err = handle(conn, addr, data, length)
		return err
	}

	return nil
}
