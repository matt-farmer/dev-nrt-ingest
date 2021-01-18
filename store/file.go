package store

import (
	"bufio"
	"os"
	"path/filepath"
	"sync"
)

// LocalFile :
type LocalFile struct {
	file  *os.File
	bw    *bufio.Writer
	count uint64
	mutex *sync.Mutex
}

// NewLocalFile :
func NewLocalFile(fileName string) (localfile *LocalFile, err error) {
	// create the output file
	if err = os.MkdirAll(filepath.Dir(fileName), os.ModePerm); err != nil {
		return nil, err
	}
	f, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}

	localfile = &LocalFile{file: f, bw: bufio.NewWriterSize(f, 65536*64), mutex: &sync.Mutex{}}
	localfile.bw.WriteString("[")
	return localfile, nil
}

// FlushClose :
func (localfile *LocalFile) FlushClose() {
	localfile.bw.WriteString("]")
	localfile.bw.Flush()
	localfile.file.Close()
}

// ----------------------------------------- //

// GenIPath :
func (localfile *LocalFile) GenIPath(path string) string {
	return ""
}

// AddVal :
func (localfile *LocalFile) AddVal(ipath string, value []byte, zip bool) {
	localfile.mutex.Lock()
	defer localfile.mutex.Unlock()

	if localfile.count > 0 {
		localfile.bw.WriteString(",\n")
	}
	localfile.bw.Write(value)
	localfile.count++

	// if atomic.AddUint64(&localfile.count, 1); localfile.count > 1 {
	// 	localfile.bw.WriteString(",\n")
	// }
	// localfile.bw.Write(value)
}

// AddValRng : [start, end)
func (localfile *LocalFile) AddValRng(ipath string, start, end int) {
}

// ----------------------------------------- //

// Get :
func (localfile *LocalFile) Get(ID string, lsPath ...string) ([]byte, bool) {
	return nil, false
}

// GetStr :
func (localfile *LocalFile) GetStr(ID string, lsPath ...string) (string, bool) {
	return "", false
}
