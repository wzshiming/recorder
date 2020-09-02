package recorder

import (
	"io"
	"time"
)

type record struct {
	duration time.Duration
	message  []byte
}

// encode encodes a record into writer.
func (r record) encode(writer io.Writer) error {
	err := putVaruint(writer, uint64(r.duration))
	if err != nil {
		return err
	}

	return putBytes(writer, r.message)
}

// decode decodes a record from reader.
func (r *record) decode(reader io.Reader) error {
	duration, err := getVaruint(reader)
	if err != nil {
		return err
	}

	message, err := getBytes(reader)
	if err != nil {
		return err
	}

	r.duration = time.Duration(duration)
	r.message = message
	return nil
}
