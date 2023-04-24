package ch04

import (
	"io"
	"net"
)

func ProxyConn(source, destination string) error {
	//1. 출발지 노드와 연결 생성
	connSource, err := net.Dial("tcp", source)
	if err != nil {
		return err
	}

	defer connSource.Close()

	//2. 목적지 노드와 연결 생성
	connDestination, err := net.Dial("tcp", destination)
	if err != nil {
		return err
	}

	defer connDestination.Close()

	//connSource에 대응하는 connDestination
	//3.connDestination으로부터 데이터를 읽고 connSource으로 데이터를 쓰기
	// 두 노드 중 연결이 끊어지면 알아서 고루틴이 종료되므로 메모리 누수 신경X
	go func() { _, _ = io.Copy(connSource, connDestination) }()

	// 4. connDestination으로 메세지를 보내는 connSource
	_, err = io.Copy(connDestination, connSource)

	return err
}
