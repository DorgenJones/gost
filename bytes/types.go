/**
*
* Copyright (c) 2018 MaoYan
* All rights reserved
* Author: dujiang02
* Date: 2019-11-20
 */
package gxbytes

import "io"

type Buffer interface {
	Read(p []byte) (n int, err error)

	ReadFrom(r io.Reader) (n int64, err error)

	Write(p []byte) (n int, err error)

	WriteTo(w io.Writer) (n int64, err error)

	WriteString(s string) (n int, err error)

	Bytes() []byte

	ReadIndex(int)

	Len() int

	Cap() int

	Init(interface{})

	Free()

	Name() string
}
