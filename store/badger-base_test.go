package store

import (
	"testing"
	"time"

	"github.com/cdutwhu/gotil/misc"
)

func TestCountDB(t *testing.T) {
	defer misc.TrackTime(time.Now())

	dbs := openBadger("../db/", "val")
	fPln(count(dbs[0]))
	closeBadger(dbs...)
}

func TestUpdateBadgerDB(t *testing.T) {
	defer misc.TrackTime(time.Now())

	dbs := openBadger("./", "test", "test1", "test2")
	defer closeBadger(dbs...)

	if updateBadgerDBs(dbs, []string{"k0", "k1", "k2"}, []string{"v0", "v1", "v2"}) != nil {
		return
	}

	if err := setBadgerBats(dbs, []string{"k0", "k1", "k2"}, []string{"v000", "v111", "v222"}); err != nil {
		fPln(err)
		return
	}
	if err := setBadgerBat(dbs[0], "k0", "k1", "VVV"); err != nil {
		fPln(err)
		return
	}

	flushBadgerBats(dbs)

	vals, err := getBadgerDBs(dbs, []string{"k0", "k1", "k2"})
	if err != nil {
		fPln(err)
	} else {
		fPln(vals)
	}

}
