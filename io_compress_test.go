package recorder

import (
	"bytes"
	"context"
	"io/ioutil"
	"testing"
	"time"
)

func Test_io_Compress(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	w := NewWriterWithCompress(context.Background(), buf)
	w.Write([]byte("Hello"))
	time.Sleep(time.Second / 10)
	w.Write([]byte(" "))
	w.Write([]byte("World"))
	time.Sleep(time.Second / 10)
	w.Write([]byte("!"))
	w.Write([]byte{})
	w.Close()
	start := time.Now()
	r, err := NewReaderWithCompress(context.Background(), buf)
	if err != nil {
		t.Fatal(err)
	}

	got, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	d := time.Since(start)
	if !bytes.Equal(got, []byte("Hello World!")) {
		t.Errorf("wrong data")
	}
	if d < time.Second/5 {
		t.Errorf("delayed failure")
	}
	r.Close()
}
