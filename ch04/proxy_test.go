package ch04

import (
	"io"
	"net"
	"sync"
	"testing"
)

// from 보내는곳 , to 받는곳 Proxy
func proxy(from io.Reader, to io.Writer) error {
	fromWriter, fromIsWriter := from.(io.Writer)
	toReader, toIsReader := to.(io.Reader)

	// 보내는곳과 받는곳이 존재한다면!
	if toIsReader && fromIsWriter {
		// 필요한 인터페이스를 구현하였으니
		// From과 to에 상응하는 프락시 생성
		go func() { _, _ = io.Copy(fromWriter, toReader) }()
	}

	_, err := io.Copy(to, from)

	return err
}

func TestProxy(t *testing.T) {
	var wg sync.WaitGroup

	//서버는 ping으로 메세지를 대기, "pong" 메세지로 응답
	// 그외는 메세지는 동일하게 클라이언트로 에코잉됨
	server, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			// 서버 수신대기
			conn, err := server.Accept()
			if err != nil {
				return
			}

			go func(c net.Conn) {
				defer c.Close()

				for {
					// 1024 버퍼 생성
					buf := make([]byte, 1024)
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							t.Error(err)
						}
						return
					}
					switch msg := string(buf[:n]); msg {
					case "ping":
						_, err = c.Write([]byte("pong"))
					default:
						_, err = c.Write(buf[:n])
					}
				}

				if err != nil {
					if err != io.EOF {
						t.Error(err)
					}

					return
				}
			}(conn)
		}
	}()

	// proxyserver는 메세지를 클라이언트 연결로부터 destinationServer (목적지)로 프록시합니다.
	// destinationServer 서버에서 온 응답메세지는 역으로 클라이언트에게 프록시합니다.
	// 1. 프록시 서버 셋업
	proxyServer, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			// 2. 프록시서버 연결 수신받아주기
			conn, err := proxyServer.Accept()
			if err != nil {
				return
			}

			go func(from net.Conn) {
				defer from.Close()
				//3. 목적지 서버 (DestinationServer) 연결
				to, err := net.Dial("tcp", server.Addr().String())
				if err != nil {
					t.Error(err)
					return
				}

				defer to.Close()

				// 4. 메세지를 프록시
				err = proxy(from, to)
				if err != nil && err != io.EOF {
					t.Error(err)
				}
			}(conn)
		}
	}()

	conn, err := net.Dial("tcp", proxyServer.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	msgs := []struct{ Message, Reply string }{
		{"ping", "pong"},
		{"ping", "pong"},
		{"echo", "echo"},
		{"ping", "pong"},
	}

	for i, m := range msgs {
		_, err := conn.Write([]byte(m.Message))
		if err != nil {
			t.Fatal(err)
		}
		buf := make([]byte, 1024)

		n, err := conn.Read(buf)
		if err != nil {
			t.Fatal(err)
		}

		actual := string(buf[:n])
		t.Logf("%q -> proxy -> %q", m.Message, actual)

		if actual != m.Reply {
			t.Errorf("%d: expected reply: %q; actual: %q",
				i, m.Reply, actual)
		}
	}
	_ = conn.Close()
	_ = proxyServer.Close()
	_ = server.Close()

	wg.Wait()
}
