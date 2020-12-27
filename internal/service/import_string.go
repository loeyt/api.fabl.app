package service

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

func extractImportString(s string) ([]byte, error) {
	if len(s) == 0 {
		return nil, errors.New("empty import string")
	}
	firstRune, s := s[0], s[1:]
	if firstRune != '0' {
		return nil, fmt.Errorf("unrecognized first character: 0x%x", firstRune)
	}
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	r, err := zlib.NewReader(bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, r)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
