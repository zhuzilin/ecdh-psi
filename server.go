package main

import (
	"fmt"
	"net"
)

type Conn struct {
	listener net.Listener
	recvConn net.Conn
	sendConn net.Conn
}

const (
	SERVER = "server"
	CLIENT = "client"
)

func NewConn(role, myAddr, peerAddr string) Conn {
	var sendConn, recvConn net.Conn
	var listener net.Listener
	var err error
	if role == SERVER {
		listener, err = net.Listen("tcp", myAddr)
		if err != nil {
			panic(err)
		}
		recvConn, err = listener.Accept()
		if err != nil {
			panic(err)
		}
		sendConn, err = net.Dial("tcp", peerAddr)
		if err != nil {
			panic(err)
		}
	} else {
		sendConn, err = net.Dial("tcp", peerAddr)
		if err != nil {
			panic(err)
		}
		listener, err = net.Listen("tcp", myAddr)
		if err != nil {
			panic(err)
		}
		recvConn, err = listener.Accept()
		if err != nil {
			panic(err)
		}
	}

	return Conn{
		listener: listener,
		recvConn: recvConn,
		sendConn: sendConn,
	}
}

func (conn *Conn) Send(val []byte) (int, error) {
	return conn.sendConn.Write(val)
}

func (conn *Conn) Receive(buf []byte) (int, error) {
	return conn.recvConn.Read(buf)
}

func (conn *Conn) SendReceiveSameLength(val []byte, buf []byte) {
	n, err := conn.Send(val)
	if err != nil {
		panic(err)
	} else if n != len(val) {
		panic(fmt.Sprintf("only send %v bytes\n", n))
	}
	n, err = conn.Receive(buf)
	if err != nil {
		panic(err)
	} else if n != len(val) {
		panic(fmt.Sprintf("only receive %v bytes\n", n))
	}
}

func (conn *Conn) Close() {
	conn.listener.Close()
	conn.recvConn.Close()
	conn.sendConn.Close()
}
