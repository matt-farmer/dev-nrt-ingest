package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	os.Remove("err-json.json")
	os.Remove("err-json.xml")	
	main()
}

func BenchmarkMain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		main()
	}
}
