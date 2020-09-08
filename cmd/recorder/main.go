package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"

	"github.com/wzshiming/notify"
	"github.com/wzshiming/recorder"
)

var (
	play   bool
	stderr bool
	stdout bool
	file   = "record.record"
)

func init() {
	flag.StringVar(&file, "f", file, "Save and load files")
	flag.BoolVar(&play, "p", play, "Play the file")
	flag.BoolVar(&stdout, "o", stdout, "Write a copy of the passed data to stdout")
	flag.BoolVar(&stderr, "e", stderr, "Write a copy of the passed data to stderr, if record")
	flag.Parse()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	notify.Once(os.Interrupt, cancel)
	var std io.Writer

	if stderr {
		std = os.Stderr
	}

	if play {
		err := doPlay(ctx, file, os.Stdout, std)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		if stdout {
			if std == nil {
				std = os.Stdout
			} else {
				std = io.MultiWriter(std, os.Stdout)
			}
		}
		err := doRecord(ctx, file, os.Stdin, std)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func doRecord(ctx context.Context, file string, in io.Reader, show io.Writer) error {
	out, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	w := recorder.NewWriterWithCompress(ctx, out)
	defer w.Close()

	var o io.Writer = w
	if show != nil {
		o = io.MultiWriter(w, show)
	}

	_, err = io.Copy(o, in)
	return err
}

func doPlay(ctx context.Context, file string, out io.Writer, show io.Writer) error {
	in, err := os.Open(file)
	if err != nil {
		return err
	}

	r, err := recorder.NewReaderWithCompress(ctx, in)
	if err != nil {
		return err
	}
	defer r.Close()

	if show != nil {
		out = io.MultiWriter(out, show)
	}

	_, err = io.Copy(out, r)
	return err
}
