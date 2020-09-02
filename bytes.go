package recorder

import (
	"io"
)

// putBytes encodes a bytes into writer.
func putBytes(writer io.Writer, x []byte) error {
	err := putVaruint(writer, uint64(len(x)))
	if err != nil {
		return err
	}
	if len(x) == 0 {
		return nil
	}
	_, err = writer.Write(x)
	return err
}

// getBytes decodes a bytes from reader.
func getBytes(reader io.Reader) ([]byte, error) {
	size, err := getVaruint(reader)
	if err != nil {
		return nil, err
	}
	if size == 0 {
		return []byte{}, nil
	}
	buf := make([]byte, size)
	i, err := io.ReadAtLeast(reader, buf, int(size))
	if err != nil {
		return nil, err
	}
	return buf[:i], nil
}
