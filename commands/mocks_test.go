package commands

import "errors"

type BrokenReader struct{}

func (r *BrokenReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("reader is broken")
}
