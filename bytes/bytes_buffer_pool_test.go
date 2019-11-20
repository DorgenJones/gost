package gxbytes

import (
	"runtime"
	"testing"
)

func TestBytesBufferPool(t *testing.T) {
	buf := GetByteBuffer(100)
	bytes := []byte{0x00, 0x01, 0x02, 0x03, 0x04}
	buf.Write(bytes)
	if buf.Len() != len(bytes) {
		t.Error("iobuffer len not match write bytes' size")
	}
	PutByteBuffer(buf)
	//buf2 := GetByteBuffer()
	// https://go-review.googlesource.com/c/go/+/162919/
	// before go 1.13, sync.Pool just reserves some objs before every gc and will be cleanup by gc.
	// after Go 1.13, maybe there are many reserved objs after gc.
	//if buf != buf2 {
	//	t.Errorf("buf pointer %p != buf2 pointer %p", buf, buf2)
	//}
}



const Size = 2048

func testbyte() []byte {
	buf := make([]byte, Size)
	for i := 0; i < Size; i++ {
		buf[i] = 1
	}
	return buf
}

func BenchmarkByteMake(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testbyte()
		if i%100 == 0 {
			runtime.GC()
		}
	}
}

// Test ByteBufferPool
var Buffer [Size]byte

func testiobufferpool() *ByteBuffer {
	b := GetByteBuffer(Size)
	b.Write(Buffer[:])
	return b
}

func testiobuffer() *ByteBuffer {
	b := NewByteBufferWithCapacity(Size)
	b.Write(Buffer[:])
	return b
}

func BenchmarkIoBufferPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := testiobufferpool()
		PutByteBuffer(buf)
		if i%100 == 0 {
			runtime.GC()
		}
	}
}

func BenchmarkIoBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testiobuffer()
		if i%100 == 0 {
			runtime.GC()
		}
	}
}

func Test_IoBufferPool(t *testing.T) {
	str := "ByteBufferPool Test"
	buffer := GetByteBuffer(len(str))
	buffer.Write([]byte(str))

	b := make([]byte, 32)
	_, err := buffer.Read(b)

	if err != nil {
		t.Fatal(err)
	}

	PutByteBuffer(buffer)

	if string(b[:len(str)]) != str {
		t.Fatal("ByteBufferPool Test Failed")
	}
	t.Log("ByteBufferPool Test Sucess")
}

func Test_IoBufferPool_Slice_Increase(t *testing.T) {
	str := "ByteBufferPool Test"
	// []byte slice increase
	buffer := GetByteBuffer(1)
	buffer.Write([]byte(str))

	b := make([]byte, 32)
	_, err := buffer.Read(b)

	if err != nil {
		t.Fatal(err)
	}

	PutByteBuffer(buffer)

	if string(b[:len(str)]) != str {
		t.Fatal("ByteBufferPool Test Slice Increase Failed")
	}
	t.Log("ByteBufferPool Test Slice Increase Sucess")
}

func Test_IoBufferPool_Alloc_Free(t *testing.T) {
	str := "ByteBufferPool Test"
	buffer := GetByteBuffer(100)
	buffer.Free()
	buffer.Init(1)
	buffer.Write([]byte(str))

	b := make([]byte, 32)
	_, err := buffer.Read(b)

	if err != nil {
		t.Fatal(err)
	}

	PutByteBuffer(buffer)

	if string(b[:len(str)]) != str {
		t.Fatal("ByteBufferPool Test Fill Free Failed")
	}
	t.Log("ByteBufferPool Test Fill Free Sucess")
}