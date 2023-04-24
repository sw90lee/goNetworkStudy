package tcp

import (
	"io"
	"net"
	"testing"
)

func TestListener(t *testing.T) {
	lisener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	defer func() { _ = lisener.Close() }()

	t.Logf("bound to %q", lisener.Addr())

	for {
		conn, err := lisener.Accept()
		if err != nil {
			t.Fatal(err)
		}

		go func(c net.Conn) {
			defer c.Close()
		}(conn)
	}
}

func TestDial(t *testing.T) {
	// TCP 를 접속할 수있는 리스너 생성
	lisener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})

	//
	go func() {
		defer func() { done <- struct{}{} }()

		for {
			conn, err := lisener.Accept()
			if err != nil {
				t.Log(err)
				return
			}

			go func(c net.Conn) {
				defer func() {
					c.Close()
					done <- struct{}{}
				}()

				buf := make([]byte, 1024)
				for {
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							t.Error(err)
						}
						return
					}
					t.Logf("received: %q", buf[:n])
				}
			}(conn)
		}
	}()

	conn, err := net.Dial("tcp", lisener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	conn.Close()
	<-done

	lisener.Close()
	<-done
}
