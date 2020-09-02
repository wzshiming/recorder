package recorder

import (
	"fmt"
	"io"
)

const (
	maxVaruintLen64 = 10
)

var (
	errOverflow = fmt.Errorf("overflow")
)

// putVaruint encodes a uint64 into writer.
func putVaruint(writer io.Writer, x uint64) error {
	var buf [maxVaruintLen64]byte
	i := 0
	for x >= 0x80 {
		buf[i] = byte(x) | 0x80
		x >>= 7
		i++
	}
	buf[i] = byte(x)
	_, err := writer.Write(buf[:i+1])
	return err
}

// getVaruint decodes a uint64 from reader.
func getVaruint(reader io.Reader) (uint64, error) {
	var buf [1]byte
	var x uint64
	var s uint
	for i := 0; ; i++ {
		_, err := reader.Read(buf[:])
		if err != nil {
			return 0, err
		}
		b := buf[0]
		if b < 0x80 {
			if i == maxVaruintLen64 || i == maxVaruintLen64-1 && b > 1 {
				return 0, errOverflow
			}
			return x | uint64(b)<<s, nil
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
}
