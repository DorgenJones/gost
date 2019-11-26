package gxbytes

import (
	"runtime"
	"testing"
)

func TestByteBufferPool(t *testing.T) {
	buf := GetByteBuffer(100)
	bytes := []byte{0x00, 0x01, 0x02, 0x03, 0x04}
	buf.Write(bytes)
	if buf.Len() != len(bytes) {
		t.Error("iobuffer len not match write bytes' size")
	}
	PutByteBuffer(buf)
}

func buildByteBuffer() Buffer {
	b := NewByteBufferWithCloneBytes(sampleBytes())
	return b
}

func Test_ByteBuffer_Grow(t *testing.T) {
	str := string(sampleBytes())
	buffer := GetByteBuffer(len(str) / 2)
	buffer.Write([]byte(str))

	b := make([]byte, Size*2)
	_, err := buffer.Read(b)
	if err != nil {
		t.Fatal(err)
	}
	if string(b[:len(str)]) != str {
		t.Fatal("ByteBufferPool Test Slice Increase Failed")
	}

	buffer.WriteString(str)
	_, err = buffer.Read(b)
	if err != nil {
		t.Fatal(err)
	}
	PutByteBuffer(buffer)
	if string(b[:len(str)]) != str {
		t.Fatal("ByteBufferPool Test Slice Increase Failed")
	}
}

func Test_Reuse(t *testing.T) {
	str := string(sampleBytes())
	buffer := GetByteBuffer(100)
	buffer.Free()
	buffer.Init(100)
	buffer.WriteString(str)
	readStr := make([]byte, len(str))
	_, err := buffer.Read(readStr)
	if err != nil {
		t.Fatal(err)
	}
	PutByteBuffer(buffer)
	if string(readStr[:len(str)]) != str {
		t.Fatal("Reuse Test Fail!!")
	}
}

func BenchmarkByteBufferPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := buildByteBuffer()
		PutByteBuffer(buf)
		if i%101 == 0 {
			runtime.GC()
		}
	}
}