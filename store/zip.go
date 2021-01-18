package store

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"strings"
)

func zipStr(strBytes []byte) []byte {
	var sb strings.Builder
	gz := gzip.NewWriter(&sb)
	if _, err := gz.Write(strBytes); err != nil {
		panic(err)
	}
	if err := gz.Flush(); err != nil {
		panic(err)
	}
	if err := gz.Close(); err != nil {
		panic(err)
	}
	str64 := base64.StdEncoding.EncodeToString([]byte(sb.String()))
	data, err := base64.StdEncoding.DecodeString(str64)
	if err != nil {
		panic(err)
	}
	return data
}

func unzipStr(zippedBytes []byte) []byte {
	r, err := gzip.NewReader(bytes.NewReader(zippedBytes))
	if err != nil {
		panic(err)
	}
	strBytes, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return strBytes
}
