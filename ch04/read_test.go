package ch04

import (
	"crypto/rand"
	"io"
	"net"
	"testing"
)

func TestReadIntoBuffer(t *testing.T) {
	payload := make([]byte, 1<<24) // 16MB
	_, err := rand.Read(payload)   // 랜덤한 페이로드 생성
	if err != nil {
		t.Fatal(err)
	}

	// 리스너 생성
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	// 고루틴으로 리스너 수신 연결 대기
	go func() {
		// 수신 연결
		conn, err := listener.Accept()
		if err != nil {
			t.Log(err)
			return
		}

		defer conn.Close()

		// 서버는 네트워크 연결로 페이로드 전체 사용
		_, err = conn.Write(payload)
		if err != nil {
			t.Error(err)
		}
	}()

	// tcp 프로토콜 , 127.0.0.1 로 연결
	conn, err := net.Dial("tcp", listener.Addr().String())

	if err != nil {
		t.Fatal(err)
	}

	// 클라이언트가 받아 드릴 수있는 데이터
	buf := make([]byte, 1<<19) // 512KB

	// 16MB payload를 512kb buffer로 받아드려야하기에 for문으로 받아드림
	for {
		n, err := conn.Read(buf)
		// 에러가 반환되거나 16MB 페이로드를 받아드릴때까지 반복
		if err != nil {
			if err != io.EOF {
				t.Error(err)
			}
			break
		}

		t.Logf("read %d bytes", n) // buf[:n]은 conn 객체에서 읽은데이터
	}
	conn.Close()
}
