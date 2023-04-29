package ch05

import (
	"bytes"
	"context"
	"net"
	"testing"
)

func TestEchoServerUDP(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	// 1.서버주소 받기
	serverAddr, err := echoServerUDP(ctx, "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	//2.ListenPacket 함수는 클라이언트와 서버 양측의 연결 객체 생성
	client, err := net.ListenPacket("udp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	defer func() { _ = client.Close() }()

	msg := []byte("ping")
	// 3. 클라이언트에게 메세지를 전송하기위해 메세지와 주소 매개변수로 받음
	_, err = client.WriteTo(msg, serverAddr)
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 1024)
	// ReadFrom 메서드에서 반환된 addr를 사용하여 에코서버가 메세지를 보냈는지 확인이 가능
	// ReadFrom 메서드를 통해 클라이언트는 즉시 메세지를 읽을 수 있음
	n, addr, err := client.ReadFrom(buf)
	if err != nil {
		t.Fatal(err)
	}

	if addr.String() != serverAddr.String() {
		t.Fatalf("received reply from %q instead of %q ", addr, serverAddr)
	}

	if !bytes.Equal(msg, buf[:n]) {
		t.Errorf("expected reply %q; actual reply %q", msg, buf[:n])
	}
}
