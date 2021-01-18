package store

import (
	"os"
	"path/filepath"

	badger "github.com/dgraph-io/badger/v2"
)

// BadgerDB :
type BadgerDB struct {
	dbPathIdx     *badger.DB
	wbPathIdx     *badger.WriteBatch
	dbIPathVal    *badger.DB
	wbIPathVal    *badger.WriteBatch
	dbIPathValRng *badger.DB
	wbIPathValRng *badger.WriteBatch
}

// NewBadgerDB :
func NewBadgerDB(dbFolder string) (badgerDB *BadgerDB, err error) {
	// remove any existing dbs
	if _, err = os.Stat(filepath.Dir(dbFolder)); os.IsExist(err) {
		if err = os.RemoveAll(filepath.Dir(dbFolder)); err != nil {
			return nil, err
		}
	}

	badgerDB = &BadgerDB{}
	for _, subfolder := range []string{"idx", "val", "rng"} {
		folder := filepath.Join(dbFolder, subfolder)
		if err = os.MkdirAll(filepath.Dir(folder), os.ModePerm); err != nil {
			return nil, err
		}

		switch subfolder {
		case "idx":
			if badgerDB.dbPathIdx, err = badger.Open(badger.DefaultOptions(folder)); err != nil {
				return nil, err
			}
			badgerDB.wbPathIdx = badgerDB.dbPathIdx.NewWriteBatch()
		case "val":
			if badgerDB.dbIPathVal, err = badger.Open(badger.DefaultOptions(folder)); err != nil {
				return nil, err
			}
			badgerDB.wbIPathVal = badgerDB.dbIPathVal.NewWriteBatch()
		case "rng":
			if badgerDB.dbIPathValRng, err = badger.Open(badger.DefaultOptions(folder)); err != nil {
				return nil, err
			}
			badgerDB.wbIPathValRng = badgerDB.dbIPathValRng.NewWriteBatch()
		}
	}

	return badgerDB, nil
}

// Close : defer
func (m *BadgerDB) Close() {
	m.wbPathIdx.Cancel()
	m.dbPathIdx.Close()
	m.wbIPathVal.Cancel()
	m.dbIPathVal.Close()
	m.wbIPathValRng.Cancel()
	m.dbIPathValRng.Close()
}

// Flush :
func (m *BadgerDB) Flush() {
	m.wbPathIdx.Flush()
	m.wbIPathVal.Flush()
	m.wbIPathValRng.Flush()
}

// ----------------------------------------- //

// GenIPath :
func (m *BadgerDB) GenIPath(path string) string {
	if len(path) <= 37 || !isUUID(path[:36]) || path[36] != '@' {
		panic("path MUST start with UUID@")
	}

	// num := -1
	// if loc := rxPathLastIdx.FindStringIndex(path); loc != nil {
	// 	s := loc[0]
	// 	num64, _ := scParseInt(path[s+1:], 10, 64)
	// 	num = int(num64)
	// 	path = path[:s]
	// }

	// m.dbPathIdx.GetSequence()

	return ""
}

// AddVal :
func (m *BadgerDB) AddVal(ipath string, value []byte, zip bool) {
	m.wbIPathVal.Set([]byte(ipath), value)
}

// AddValRng : [start, end)
func (m *BadgerDB) AddValRng(ipath string, start, end int) {
}

// ----------------------------------------- //

// Get :
func (m *BadgerDB) Get(ID string, lsPath ...string) ([]byte, bool) {
	return nil, false
}

// GetStr :
func (m *BadgerDB) GetStr(ID string, lsPath ...string) (string, bool) {
	return "", false
}
