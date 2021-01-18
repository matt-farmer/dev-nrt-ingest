package store

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/cdutwhu/debog/fn"
	"github.com/google/uuid"
)

const (
	sep, preatt = "~", "-"
)

// Path
var (
	rxPathIdx     = regexp.MustCompile(`#\d+`)
	rxPathLastIdx = regexp.MustCompile(`#\d+$`)
)

var (
	fPln       = fmt.Println
	fSf        = fmt.Sprintf
	fEf        = fmt.Errorf
	sJoin      = strings.Join
	scParseInt = strconv.ParseInt

	failOnErr       = fn.FailOnErr
	failP1OnErrWhen = fn.FailP1OnErrWhen

	isUUID = func(u string) bool {
		_, err := uuid.Parse(u)
		return err == nil
	}
)

var (
	// Explicitly indicate which PATH is LIST
	mList = make(map[string]struct{})

	// if at least one 'AddVal'
	zipping = false
)

// SetLISTPath :
func SetLISTPath(listpaths ...string) {
	for _, p := range listpaths {
		mList[p] = struct{}{}
	}
}
