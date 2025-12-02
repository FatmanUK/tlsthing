package main

import (
	"io"
	"fmt"
	"net"
)

var MAXBUFLEN int = 2048

func ReadTcpBytes(conn net.Conn, logs chan Log) ([]byte, error) {
	ident := conn.RemoteAddr().String()
	buffer := make([]byte, MAXBUFLEN)
	n, err := conn.Read(buffer)
	if err != nil {
		if err == io.EOF {
			msg := fmt.Sprintf("Closing %s.", ident)
			logs <- Info(msg)
		} else {
			logs <- Error(err.Error())
		}
		return nil, err
	}
	logs <- Debug(fmt.Sprintf("Received %d bytes.", n))
	return buffer[:n], nil
}

func WriteTcpBytes(conn net.Conn, logs chan Log, buf []byte) error {
	ident := conn.RemoteAddr().String()
	msg := fmt.Sprintf("Closing %s.", ident)
	logs <- Info(msg)
	n, err := conn.Write(buf)
	logs <- Debug(fmt.Sprintf("Sent %d bytes.", n))
	return err
}

func GetDirectorConnection(
		hostPort string,
		logs chan Log) (net.Conn, error) {
	var conn net.Conn
	var err error
	count := 0
	dirConnectMsg := "Connecting director (attempt %d)..."
	for ok := true; ok; ok = (err != nil && count < 9) {
		DoublingSleep(count, 2)
		logs <- Info(fmt.Sprintf(dirConnectMsg, count + 1))
		conn, err = net.Dial("tcp", hostPort)
		if err == nil {
			break
		}
		count++
	}
	if err != nil {
		logs <- Error("Connection failed.")
		return nil, err
	}
	logs <- Info("Connection established.")
	return conn, nil
}

func ThreadAcceptHandler(
		conn net.Conn,
		logs chan Log,
		data interface{}) {
	defer conn.Close()
	ident := conn.RemoteAddr().String()
	logs <- Info(fmt.Sprintf("Connected %s.", ident))
	for {
		buffer, err := ReadTcpBytes(conn, logs)
		if err != nil {
			break
		}
		DebugOutputByteArray(buffer, logs)
		msgs, err := DecodeNetworkMessages(buffer, logs)
		if err != nil {
			continue
		}
		err = ProcessMessages(conn, msgs, logs, data)
		if err != nil {
			continue
		}
	}
}
