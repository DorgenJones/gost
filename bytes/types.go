/**
*
* Copyright (c) 2018 MaoYan
* All rights reserved
* Author: dujiang02
* Date: 2019-11-20
 */
package gxbytes

import (
	"io"
)

type Buffer interface {
	// Read reads the next len(p) bytes from the buffer or until the buffer
	// is drained. The return value n is the number of bytes read. If the
	// buffer has no data to return, err is io.EOF (unless len(p) is zero);
	// otherwise it is nil.
	Read(p []byte) (n int, err error)

	// Grow grows the buffer's capacity, if necessary, to guarantee space for
	// another n bytes. After Grow(n), at least n bytes can be written to the
	// buffer without another allocation.
	// If n is negative, Grow will panic.
	// If the buffer can't grow it will panic with ErrTooLarge.
	Grow(n int)

	// ReadFrom reads data from r until EOF and appends it to the buffer, growing
	// the buffer as needed. The return value n is the number of bytes read. Any
	// error encountered during the read is also returned. If the
	// buffer becomes too large, ReadFrom will panic with ErrTooLarge.
	ReadFrom(r io.Reader) (n int64, err error)

	// Write appends the contents of p to the buffer, growing the buffer as
	// needed. The return value n is the length of p; err is always nil. If the
	// buffer becomes too large, Write will panic with ErrTooLarge.
	Write(p []byte) (n int, err error)

	// WriteTo writes data to w until the buffer is drained or an error occurs.
	// The return value n is the number of bytes written; it always fits into an
	// int, but it is int64 to match the io.WriterTo interface. Any error
	// encountered during the write is also returned.
	WriteTo(w io.Writer) (n int64, err error)

	// WriteString appends the contents of s to the buffer, growing the buffer as
	// needed. The return value n is the length of s; err is always nil. If the
	// buffer becomes too large, WriteString will panic with ErrTooLarge.
	WriteString(s string) (n int, err error)

	// Bytes returns a slice of length b.Len() holding the unread portion of the buffer.
	// The slice is valid for use only until the next buffer modification (that is,
	// only until the next call to a method like Read, Write, Reset, or Truncate).
	// The slice aliases the buffer content at least until the next buffer modification,
	// so immediate changes to the slice will affect the result of future reads.
	Bytes() []byte

	//shift n byte
	Shift(int)

	// Len returns the number of bytes of the unread portion of the buffer;
	// b.Len() == len(b.Bytes()).
	Len() int

	// Cap returns the capacity of the buffer's underlying byte slice, that is, the
	// total space allocated for the buffer's data.
	Cap() int

	// init
	Init(interface{})

	// free resource
	Free()

	// name
	Name() string
}