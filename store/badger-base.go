package store

import badger "github.com/dgraph-io/badger/v2"

var (
	mDbWb = make(map[*badger.DB]*badger.WriteBatch)
)

func openBadger(dbPath string, lsDBName ...string) []*badger.DB {
	defer func() { fPln("openBadger done...") }()
	dbs := make([]*badger.DB, len(lsDBName))
	for i, name := range lsDBName {
		db, err := badger.Open(badger.DefaultOptions(dbPath + name))
		failOnErr("%v", err)
		dbs[i] = db
		mDbWb[db] = db.NewWriteBatch()
	}
	return dbs
}

func closeBadger(dbs ...*badger.DB) error {
	defer func() { fPln("closeBadger done...") }()
	for _, db := range dbs {
		if db != nil {
			mDbWb[db].Cancel()
			mDbWb[db] = nil
			if err := db.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

// count :
func count(db *badger.DB) (cnt int) {
	opt := badger.DefaultIteratorOptions
	db.View(func(txn *badger.Txn) error {
		itr := txn.NewIterator(opt)
		defer itr.Close()
		for itr.Rewind(); itr.Valid(); itr.Next() {
			item := itr.Item()
			item.Value(func(v []byte) error {
				cnt++
				return nil
			})
		}
		return nil
	})
	return cnt
}

func setBadgerBat(db *badger.DB, ksvs ...string) error {
	if len(ksvs)%2 != 0 {
		return fEf("each of ksvs key MUST have one value; left ksvs are keys, right ksvs are values")
	}
	l := len(ksvs) / 2
	wb := mDbWb[db]
	for i := 0; i < l; i++ {
		key, val := ksvs[i], ksvs[l+i]
		if err := wb.Set([]byte(key), []byte(val)); err != nil {
			return err
		}
	}
	return nil
}

// if 'lsValues' is not provided, delete keys
func setBadgerBats(dbs []*badger.DB, keys []string, lsValues ...[]string) (err error) {
	if len(dbs) != len(keys) {
		return fEf("db count MUST be equal to key's")
	}
	if len(lsValues) > 1 {
		return fEf("lsValues at most provide one []string")
	}
	if len(lsValues) == 1 && len(keys) != len(lsValues[0]) {
		return fEf("key count MUST be equal to value's")
	}

	for i, db := range dbs {
		wb := mDbWb[db]
		if key := []byte(keys[i]); len(lsValues) == 0 { // delete
			err = wb.Delete(key)
		} else { //                                        set
			val := []byte(lsValues[0][i])
			err = wb.Set(key, val)
		}
		if err != nil {
			return err
		}
	}
	return err
}

func flushBadgerBats(dbs []*badger.DB) error {
	for _, db := range dbs {
		if err := mDbWb[db].Flush(); err != nil {
			return err
		}
	}
	return nil
}

// if 'lsValues' is not provided, delete keys
func updateBadgerDBs(dbs []*badger.DB, keys []string, lsValues ...[]string) (err error) {
	if len(dbs) != len(keys) {
		return fEf("db count MUST be equal to key's")
	}
	if len(lsValues) > 1 {
		return fEf("lsValues at most provide one []string")
	}
	if len(lsValues) == 1 && len(keys) != len(lsValues[0]) {
		return fEf("key count MUST be equal to value's")
	}

	commit := func(lsTxn ...*badger.Txn) error {
		for _, tx := range lsTxn {
			if tx != nil {
				if err := tx.Commit(); err != nil {
					return err
				}
			}
		}
		return nil
	}

	lsTxn := []*badger.Txn{}
	for i, db := range dbs {
		txn := db.NewTransaction(true)
		if key := []byte(keys[i]); len(lsValues) == 0 { // delete
			err = txn.Delete(key)
		} else { //                                        set
			val := []byte(lsValues[0][i])
			err = txn.Set(key, val)
		}
		if err != nil {
			break
		}
		lsTxn = append(lsTxn, txn) // only keep non-error txn to discard
	}
	defer func() {
		for _, txn := range lsTxn {
			txn.Discard()
		}
	}()
	if err == nil {
		return commit(lsTxn...)
	}
	return err
}

func getBadgerDBs(dbs []*badger.DB, keys []string) (values []string, err error) {
	if len(dbs) != len(keys) {
		return nil, fEf("db count MUST be equal to key's")
	}

	lsTxn := []*badger.Txn{}
	for i, db := range dbs {
		txn := db.NewTransaction(true)
		lsTxn = append(lsTxn, txn)
		switch item, e := txn.Get([]byte(keys[i])); e {
		case nil:
			item.Value(func(v []byte) error {
				values = append(values, string(v))
				return nil
			})
		case badger.ErrKeyNotFound:
			values = append(values, "")
		default:
			failOnErr("%v", e)
		}
	}
	defer func() {
		for _, txn := range lsTxn {
			txn.Discard()
		}
	}()
	return
}

func getBadgerDB(db *badger.DB, keys []string) (values []string, err error) {
	txn := db.NewTransaction(true)
	defer txn.Discard()
	for _, key := range keys {
		switch item, e := txn.Get([]byte(key)); e {
		case nil:
			item.Value(func(v []byte) error {
				values = append(values, string(v))
				return nil
			})
		case badger.ErrKeyNotFound:
			values = append(values, "")
		default:
			failOnErr("%v", e)
		}
	}
	return
}
