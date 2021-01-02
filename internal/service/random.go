package service

import "math/rand"

type globalRandReader struct{}

func (globalRandReader) Read(p []byte) (n int, err error) {
	return rand.Read(p)
}
