package recorder

import (
	"bytes"
	"reflect"
	"testing"
)

func Test_record(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	want := record{
		duration: 100,
		message:  []byte{1, 2, 3, 4},
	}
	err := want.encode(buf)
	if err != nil {
		t.Fatal(err)
	}
	err = want.encode(buf)
	if err != nil {
		t.Fatal(err)
	}

	got := record{}
	err = got.decode(buf)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want is not equal got")
	}

	err = got.decode(buf)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want is not equal got")
	}
}
