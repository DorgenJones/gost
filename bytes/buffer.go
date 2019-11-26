// Copy from types.buffer (go@1.12.4) and change
// 1. clean useless code
// 2. change the process of ReadFrom

package gxbytes

import (
	"errors"
	"io"
)

type ByteBuffer struct {
	buf      []byte
	off      int
	lastRead readOp
	copy     *[]byte
}

type readOp int8

const (
	opRead      readOp = -1 // Any other read operation.
	opInvalid   readOp = 0  // Non-read operation.
	opReadRune1 readOp = 1  // Read rune of size 1.
	opReadRune2 readOp = 2  // Read rune of size 2.
	opReadRune3 readOp = 3  // Read rune of size 3.
	opReadRune4 readOp = 4  // Read rune of size 4.
)

var ErrTooLarge = errors.New("bytes.ByteBuffer: too large")
var errNegativeRead = errors.New("bytes.ByteBuffer: reader returned negative count from Read")

const maxInt = int(^uint(0) >> 1)

func (b *ByteBuffer) Bytes() []byte { return b.buf[b.off:] }

func (b *ByteBuffer) String() string {
	if b == nil {
		return "<nil>"
	}
	return string(b.buf[b.off:])
}

func (b *ByteBuffer) empty() bool { return len(b.buf) <= b.off }

func (b *ByteBuffer) Len() int { return len(b.buf) - b.off }

func (b *ByteBuffer) Cap() int { return cap(b.buf) }

func (b *ByteBuffer) Truncate(n int) {
	if n == 0 {
		b.Reset()
		return
	}
	b.lastRead = opInvalid
	if n < 0 || n > b.Len() {
		panic("bytes.ByteBuffer: truncation out of range")
	}
	b.buf = b.buf[:b.off+n]
}

func (b *ByteBuffer) Reset() {
	b.buf = b.buf[:0]
	b.off = 0
	b.lastRead = opInvalid
}

func (b *ByteBuffer) Free() {
	b.Reset()
	if b.copy != nil {
		PutBytes(b.copy)
	}
	b.off = 0
}

func (b *ByteBuffer) tryGrowByReslice(n int) (int, bool) {
	if l := len(b.buf); n <= cap(b.buf)-l {
		b.buf = b.buf[:l+n]
		return l, true
	}
	return 0, false
}

func (b *ByteBuffer) Name() string {
	return "ByteBuffer"
}

func (b *ByteBuffer) grow(n int) int {
	m := b.Len()
	// If buffer is empty, reset to recover space.
	if m == 0 && b.off != 0 {
		b.Reset()
	}
	// Try to translation
	if free := cap(b.buf) - len(b.buf) + b.off; free > n {
		newBuf := b.buf
		copy(newBuf, b.buf[b.off:])
		b.buf = newBuf[:len(b.buf)-b.off]
		b.off = 0
	}
	// Try to grow by means of a reslice.
	if i, ok := b.tryGrowByReslice(n); ok {
		return i
	}
	c := cap(b.buf)
	if n <= c/2-m {
		// We can slide things down instead of allocating a new
		// slice. We only need m+n <= c to slide, but
		// we instead let capacity get twice as large so we
		// don't spend all our time copying.
		copy(b.buf, b.buf[b.off:])
	} else if c > maxInt-c-n {
		panic(ErrTooLarge)
	} else {
		// Not enough space anywhere, we need to allocate.
		buf := *makeSlice(2*c + n)
		copy(buf, b.buf[b.off:])
		PutBytes(b.copy)
		b.buf = buf
		b.copy = &b.buf
	}
	// Restore b.off and len(b.buf).
	b.off = 0
	b.buf = b.buf[:m+n]
	return m
}

func (b *ByteBuffer) Grow(n int) {
	if n < 0 {
		panic("bytes.ByteBuffer.Grow: negative count")
	}
	m := b.grow(n)
	b.buf = b.buf[:m]
}

func (b *ByteBuffer) Write(p []byte) (n int, err error) {
	b.lastRead = opInvalid
	m, ok := b.tryGrowByReslice(len(p))
	if !ok {
		m = b.grow(len(p))
	}
	return copy(b.buf[m:], p), nil
}

func (b *ByteBuffer) WriteString(s string) (n int, err error) {
	b.lastRead = opInvalid
	m, ok := b.tryGrowByReslice(len(s))
	if !ok {
		m = b.grow(len(s))
	}
	return copy(b.buf[m:], s), nil
}

const MinRead = 1 << 9

const MaxRead = MinRead << 8

const DefaultSize = MinRead >> 2

// change:
//	1. maxRead limit
//	2. not enough bytes then return
func (b *ByteBuffer) ReadFrom(r io.Reader) (n int64, err error) {
	b.lastRead = opInvalid
	for {
		i := b.grow(MinRead)
		b.buf = b.buf[:i]
		should := cap(b.buf) - i
		m, e := r.Read(b.buf[i:cap(b.buf)])
		if m < 0 {
			panic(errNegativeRead)
		}

		b.buf = b.buf[:i+m]
		n += int64(m)
		if e != nil {
			return n, e
		}
		if should != m {
			return n, nil
		}
		if n > MaxRead {
			return n, nil
		}
	}
}

func makeSlice(n int) *[]byte {
	return GetBytes(n)
}

func (b *ByteBuffer) WriteTo(w io.Writer) (n int64, err error) {
	b.lastRead = opInvalid
	if nBytes := b.Len(); nBytes > 0 {
		m, e := w.Write(b.buf[b.off:])
		if m > nBytes {
			panic("bytes.ByteBuffer.WriteTo: invalid Write count")
		}
		b.off += m
		n = int64(m)
		if e != nil {
			return n, e
		}
		// all bytes should have been written, by definition of
		// Write method in io.Writer
		if m != nBytes {
			return n, io.ErrShortWrite
		}
	}
	b.Reset()
	return n, nil
}

func (b *ByteBuffer) Read(p []byte) (n int, err error) {
	b.lastRead = opInvalid
	if b.empty() {
		// ByteBuffer is empty, reset to recover space.
		b.Reset()
		if len(p) == 0 {
			return 0, nil
		}
		return 0, io.EOF
	}
	n = copy(p, b.buf[b.off:])
	b.off += n
	if n > 0 {
		b.lastRead = opRead
	}
	return n, nil
}

func (b *ByteBuffer) Shift(n int) {
	if b.off+n > len(b.buf) {
		return
	}
	b.off += n
	if n > 0 {
		b.lastRead = opRead
	}
}

func (b *ByteBuffer) Init(param interface{}) {
	size, ok := param.(int)
	if ok {
		if size <= 0 {
			size = MinRead
		}
		b.buf = *makeSlice(size)
		b.copy = &b.buf
		b.buf = b.buf[:0]
		b.off = 0
	}
}

func NewByteBufferWithCapacity(capacity int) *ByteBuffer {
	buffer := &ByteBuffer{
	}
	buffer.Init(capacity)
	return buffer
}

func NewByteBufferString(s string) *ByteBuffer {
	buf := []byte(s)
	return &ByteBuffer{
		buf:  buf,
		off:  0,
		copy: &buf,
	}
}

func NewByteBuffer(bytes []byte) *ByteBuffer {
	return &ByteBuffer{
		buf:  bytes,
		copy: &bytes,
		off:  0,
	}
}

func NewByteBufferWithCloneBytes(bytes []byte) *ByteBuffer {
	buffer := NewByteBufferWithCapacity(len(bytes))
	buffer.Write(bytes)
	return buffer
}
