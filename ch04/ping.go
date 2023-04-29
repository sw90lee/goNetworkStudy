package ch04

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

// 1. ping 커맨드가 제공하는 기능의 일부를 흉내
var (
	count    = flag.Int("c", 3, "number of pings: <= 0 means forever")
	interval = flag.Duration("i", time.Second, "interval between pings")
	timeout  = flag.Duration("W", 5*time.Second, "time to wait for a reply")
)

func int() {
	flag.Usage = func() {
		fmt.Println("Usage: %s [options] host:port\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {

	flag.Parse()

	if flag.NFlag() != 1 {
		fmt.Print("host: port is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	target := flag.Arg(0)
	fmt.Println("PING", target)

	if *count <= 0 {
		fmt.Println("CTRL+C to stop.")
	}

	msg := 0

	for (*count <= 0) || (msg < *count) {
		msg++
		fmt.Print(msg, " ")

		start := time.Now()
		//1. tcp 연결 수립 시도 -> 응답 하지않을경우 timeout 설정
		c, err := net.DialTimeout("tcp", target, *timeout)
		//2. tcp 핸드 세이크를 마치는데 걸리는 시간을 추적
		//( 출발지 호스트와 원격 호스트간 ping 도달하는시간 )
		dur := time.Since(start)

		if err != nil {
			fmt.Printf("fail in %s: %v\n", dur, err)
			//3. 일시적 에러가 아니면 종료
			if nErr, ok := err.(net.Error); !ok || !nErr.Temporary() {
				os.Exit(1)
			}
		} else {
			_ = c.Close()
			fmt.Println(dur)
		}

		time.Sleep(*interval)
	}
}
