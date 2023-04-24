package ch04

import (
	"io"
	"log"
	"net"
	"os"
)

// Monitor 구조체는 네트워크 트래픽을 로깅하기 위한 log.Logger를 임베딩합니다.
type Monitor struct {
	*log.Logger
}

// Write 메서드는 io.Writer 인터페이스를 구현
func (m *Monitor) Write(p []byte) (int, error) {
	return len(p), m.Output(2, string(p))
}

func ExampleMonitor() {
	//2. stdout 표준 출력으로 데이터를 쓰는 Monitor 구조체 인스턴스 생성
	monitor := &Monitor{Logger: log.New(os.Stdout, "monitor: ", 0)}

	// 서버생성
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		monitor.Fatal(err)
	}

	done := make(chan struct{})

	go func() {
		defer close(done)

		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		buf := make([]byte, 1024)
		//3. monitor 인스턴스 변수와 함께 연결 객체 사용
		r := io.TeeReader(conn, monitor)

		// io.reader로 데이터를 읽고 모니티에 출력 후
		// 함수를 호출자에게 전달
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			monitor.Println(err)
			return
		}

		//4. 서버에 출력결과를 생성한 io.MultiWriter를 이용하여
		// 네트워크 연결과 모니터링에 로깅
		w := io.MultiWriter(conn, monitor)

		_, err = w.Write(buf[:n]) // 메세지 에코잉
		if err != nil && err != io.EOF {
			monitor.Println(err)
			return
		}

	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		monitor.Fatal(err)
	}

	_, err = conn.Write([]byte("Test\n"))
	if err != nil {
		monitor.Fatal(err)
	}

	_ = conn.Close()
	<-done

}
