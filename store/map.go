package store

import (
	"sync"
)

// SyncMap : implement CRUD
type SyncMap struct {
	mPathIdx     sync.Map
	mIPathVal    sync.Map
	mIPathValRng sync.Map
}

// NewSyncMap :
func NewSyncMap() *SyncMap {
	return &SyncMap{}
}

// ----------------------------------------- //

// GenIPath :
func (m *SyncMap) GenIPath(path string) string {
	if len(path) <= 37 || !isUUID(path[:36]) || path[36] != '@' {
		panic("path MUST start with UUID@")
	}

	num := -1
	if loc := rxPathLastIdx.FindStringIndex(path); loc != nil {
		s := loc[0]
		num64, _ := scParseInt(path[s+1:], 10, 64)
		num = int(num64)
		path = path[:s]
	}

	if _, ok := m.mPathIdx.Load(path); !ok {
		m.mPathIdx.Store(path, num)
	}
	numVal, _ := m.mPathIdx.Load(path)
	m.mPathIdx.Store(path, numVal.(int)+1)
	idxVal, _ := m.mPathIdx.Load(path)
	idx := idxVal.(int)

	ipath := path

	if idx > 0 {
		ipath = fSf("%s#%d", path, idx)
	} else {
		// ENFORCE LIST
		p := rxPathIdx.ReplaceAllString(path[37:], "")
		if _, ok := mList[p]; ok {
			ipath += "#0"
		}
	}
	return ipath
}

// AddVal :
func (m *SyncMap) AddVal(ipath string, value []byte, zip bool) {
	if !zip {
		m.mIPathVal.Store(ipath, value)
	} else {
		// zip value & save
		m.mIPathVal.Store(ipath, zipStr(value))
	}
}

// AddValRng : [start, end)
func (m *SyncMap) AddValRng(ipath string, start, end int) {
	m.mIPathValRng.Store(ipath, [2]int{start, end})
}

// ----------------------------------------- //

// Get :
func (m *SyncMap) Get(ID string, lsPath ...string) ([]byte, bool) {
	xmlVal, ok := m.mIPathVal.Load(ID)
	if !ok {
		return nil, false
	}

	path := ID + "@" + sJoin(lsPath, sep)

	// CPLX
	if v, ok := m.mIPathValRng.Load(path); ok {
		rng := v.([2]int)
		xml := xmlVal.([]byte)
		// unzip xml ............................ UNZIP (../parse.go ZIP)
		// xml = unzipStr(xml)
		return xml[rng[0]:rng[1]], ok
	}

	// SIMPLE
	if v, ok := m.mIPathVal.Load(path); ok {
		return v.([]byte), ok
	}

	// ATTRIBUTE
	lastIdx := len(lsPath) - 1
	path = ID + "@" + sJoin(lsPath[:lastIdx], sep) + sep + preatt + lsPath[lastIdx]
	if v, ok := m.mIPathVal.Load(path); ok {
		return v.([]byte), ok
	}

	return nil, false
}

// GetStr :
func (m *SyncMap) GetStr(ID string, lsPath ...string) (string, bool) {
	data, ok := m.Get(ID, lsPath...)
	return string(data), ok
}
