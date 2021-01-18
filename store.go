package main

// IIngest :
type IIngest interface {
	// ZipCplx(zip bool)
	GenIPath(path string) string
	AddVal(ipath string, value []byte, zip bool)
	AddValRng(ipath string, start, end int)
	// ------------------ //
	GetStr(id string, lsPath ...string) (string, bool)
}
