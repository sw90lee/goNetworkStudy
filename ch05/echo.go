package ch05

import (
	"context"
	"fmt"
	"net"
)

// 1. context를 매개변수, 호스트:포트 형식으로 구성된 문자열주소로 받고
// 에코서버에 메세지 전송하는 함수
func echoServerUDP(ctx context.Context, addr string) (net.Addr, error) {
	//2. 에코서버에 UDP 연결
	s, err := net.ListenPacket("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("binding to udp %s: %w", addr, err)
	}
	//3. 비동기적 메세지 에코잉
	go func() {
		go func() {
			//4. 채널 블록킹 context가 취소되면 블로킹 해제되고 서버가 닫힘
			//3번의 고루틴도 종료
			<-ctx.Done()
			_ = s.Close()
		}()

		buf := make([]byte, 1024)

		for {
			// UDP 연결로부터 데이터을 읽기위해 바이트슬라이스를 매개변수로 전달
			n, clientAddr, err := s.ReadFrom(buf) // 클라이언트에서 서버로
			if err != nil {
				return
			}

			// UDP 패킷을 전송하기위해 바이트 슬라이스와 목적지주소를 연결
			// WriteTo 메서드의 매개변수로 전달
			_, err = s.WriteTo(buf[:n], clientAddr) // 서버에서 클라이언ㄴ트로
			if err != nil {
				return
			}
		}
	}()

	return s.LocalAddr(), nil
}
