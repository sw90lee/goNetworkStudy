package ch04

import (
	"bufio"
	"net"
	"reflect"
	"testing"
)

const payload = "The bigger the interface, the weaker the abstraction."

func TestScanner(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		// 수신 연결 대기
		conn, err := listener.Accept()
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()

		_, err = conn.Write([]byte(payload))
		if err != nil {
			t.Error(err)
		}
	}()

	//  클라이언트의 수신 연결
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	defer conn.Close()

	// 서버에서 문자열을 읽어들일 bufio.scannner 생성
	scanner := bufio.NewScanner(conn)
	// 공백이나 마침표등의 단어 경계를 구분하는 ScanWords를 사용
	scanner.Split(bufio.ScanWords)

	var words []string

	// 스캐너가 계속해서 데이터를 읽음 ,io.EOF or ERROR 반환할때까지 돔
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	err = scanner.Err()
	if err != nil {
		t.Error(err)
	}

	// 예상 되는 단어
	expected := []string{"The", "bigger", "the", "interface,", "the", "weaker", "the", "abstraction."}
	// 스캐너가 읽은 단어와 예상되는 단어 TEST
	if !reflect.DeepEqual(words, expected) {
		t.Fatal("inaccurate scanned word list")
	}
	t.Logf("Scanned words: %#v", words)
}
