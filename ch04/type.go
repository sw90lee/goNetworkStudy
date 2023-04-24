package ch04

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	BinaryType uint8 = iota + 1
	StringType

	MaxPayloadSize uint32 = 10 << 20 // 10MB
)

var ErrMaxPayloadSize = errors.New("maximum payload size exceeded")

type Payload interface {
	fmt.Stringer
	io.ReaderFrom
	io.WriterTo
	Bytes() []byte
}

type Binary []byte

func (m Binary) Bytes() []byte  { return m }
func (m Binary) String() string { return string(m) }

func (m Binary) WriteTo(w io.Writer) (int64, error) {
	err := binary.Write(w, binary.BigEndian, BinaryType) // 1바이트 타입
	if err != nil {
		return 0, err
	}

	var n int64 = 1

	err = binary.Write(w, binary.BigEndian, uint32(len(m))) //4바이트크기
	if err != nil {
		return n, err
	}

	n += 4

	o, err := w.Write(m) // 페이로드

	return n + int64(o), err
}

func (m *Binary) ReadFrom(r io.Reader) (int64, error) {
	var typ uint8
	err := binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		return 0, err
	}

	var n int64 = 1
	if typ != BinaryType { // typ uin8 1바이트 변수를 r 로읽어 바이너리 타입이 맞는지 확인
		return n, errors.New("invalid Binary")
	}

	var size uint32
	err = binary.Read(r, binary.BigEndian, &size) // 4바이트 크기를 읽어드림
	if err != nil {
		return n, err
	}

	n += 4
	if size > MaxPayloadSize { // 최대 크기의 페이로드 사이즈를 주었는지 확인
		return n, err
	}
	*m = make([]byte, size) // size 변수값 Binary 인스턴스의 크기로 새로운 바이트 슬라이스를 할당
	o, err := r.Read(*m)    // binary 인스턴스 바이트 슬라이스를 읽음

	return n + int64(o), err
}

type String string

func (m String) Bytes() []byte  { return []byte(m) }
func (m String) String() string { return string(m) }

func (m String) WriteTo(w io.Writer) (int64, error) {
	err := binary.Write(w, binary.BigEndian, StringType)
	if err != nil {
		return 0, err
	}

	var n int64 = 1

	err = binary.Write(w, binary.BigEndian, uint32(n)) // 4바이트 크기
	if err != nil {
		return n, err
	}

	n += 4

	o, err := w.Write([]byte(m))

	return n + int64(o), err
}

func (m *String) ReadFrom(r io.Reader) (int64, error) {
	var typ uint8
	err := binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		return 0, err
	}

	var n int64 = 1
	if typ != StringType {
		return n, errors.New("invalid String")
	}

	var size uint32
	err = binary.Read(r, binary.BigEndian, &size)
	if err != nil {
		return n, err
	}

	n += 4

	buf := make([]byte, size)
	o, err := r.Read(buf)
	if err != nil {
		return n, err
	}

	*m = String(buf)

	return n + int64(o), nil
}

// 1.
func decode(r io.Reader) (Payload, error) {
	var typ uint8
	//2. payload 변수를 생성
	err := binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		return nil, err
	}

	//3. 디코딩된 타입에 저
	var payload Payload

	//4. pay로드 변수에 대한 상수타입 할당
	switch typ {
	case BinaryType:
		payload = new(Binary)
	case StringType:
		payload = new(String)
	default:
		return nil, errors.New("unknown Type")
	}

	_, err = payload.ReadFrom(
		//5.
		io.MultiReader(bytes.NewReader([]byte{typ}), r))
	if err != nil {
		return nil, err
	}

	return payload, nil
}
