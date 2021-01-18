package store

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/cdutwhu/gotil/misc"
)

func TestZipStr(t *testing.T) {
	misc.TrackTime(time.Now())
	bytes, _ := ioutil.ReadFile("../sif.xml")
	str := string(bytes)
	fPln(len(str))
	data := zipStr(bytes)
	fPln(float64(len(str)) / float64(len(data)))
	str = string(unzipStr(data))
	fPln(len(str))
}
