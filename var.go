package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gosuri/uiprogress"
)

var (
	fPln        = fmt.Println
	fPf         = fmt.Printf
	sTrimLeft   = strings.TrimLeft
	sTrimRight  = strings.TrimRight
	sTrimSuffix = strings.TrimSuffix
	sHasSuffix  = strings.HasSuffix
	sContains   = strings.Contains

	dbgPln = func(do bool, a ...interface{}) (n int, err error) {
		if do {
			return fPln(a...)
		}
		return 0, nil
	}

	// fileExists checks if a file exists and is not a directory before we
	// try using it to prevent further errors.
	fileExists = func(filename string) bool {
		info, err := os.Stat(filename)
		if os.IsNotExist(err) {
			return false
		}
		return !info.IsDir()
	}
)

const (
	attrPrefix   = ""
	contAttrName = "value"
	attrNameOfID = "RefId"
)

var (
	validate = false
	// CheckData : validate xml & json in parsing
	CheckData = func(check bool) {
		validate = check
	}
)

var (
	uip      *uiprogress.Progress
	bar      *uiprogress.Bar
	procsize uint64
	probar   bool
)
