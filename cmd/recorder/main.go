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
	play bool
	file = "record.record"
)

func init() {
	flag.BoolVar(&play, "p", false, "Play the file")
	flag.StringVar(&file, "f", file, "Save and load files")
	flag.Parse()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	notify.Once(os.Interrupt, cancel)
	if play {
		err := doPlay(ctx)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := doRecord(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func doRecord(ctx context.Context) error {
	out, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	w := recorder.NewWriterWithCompress(ctx, out)
	defer w.Close()

	_, err = io.Copy(io.MultiWriter(os.Stdout, w), os.Stdin)
	if err != nil {
		return err
	}
	return nil
}

func doPlay(ctx context.Context) error {
	in, err := os.Open(file)
	if err != nil {
		return err
	}

	r, err := recorder.NewReaderWithCompress(ctx, in)
	if err != nil {
		return err
	}
	defer r.Close()

	_, err = io.Copy(os.Stdout, r)
	if err != nil {
		return err
	}
	return nil
}
