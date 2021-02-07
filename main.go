package main

import (
	"flag"
	"path/filepath"
	"time"

	"github.com/nsip/dev-nrt-ingest/store"
	"github.com/cdutwhu/gotil/misc"
	xt "github.com/cdutwhu/xml-tool"
)

func main() {
	defer misc.TrackTime(time.Now())

	xmlPathPtr := flag.String("input", "./rrd.xml", "path of input xml file")
	asyncPtr := flag.Bool("async", true, "async process")
	dataChkPtr := flag.Bool("check", false, "validate json and xml data")
	storeTypePtr := flag.String("store", "map", "store type [map, kvdb/badger, file]")
	probarPtr := flag.Bool("bar", true, "show progress bar")
	flag.Parse()

	if !fileExists(*xmlPathPtr) {
		fPln("xml file is not exist, file path [-input] is needed")
		return
	}
	xmlbasename := sTrimSuffix(filepath.Base(*xmlPathPtr), ".xml") // for dumped json

	CheckData(*dataChkPtr)
	probar = *probarPtr

	// ------------------------------------- //

	var ingest IIngest
	var err error

	switch *storeTypePtr {
	case "map":
		ingest = store.NewSyncMap()

	case "kvdb", "badger":
		ingest, err = store.NewBadgerDB("./db")
		if err != nil {
			panic(err)
		}
		defer ingest.(*store.BadgerDB).Close()
		defer ingest.(*store.BadgerDB).Flush()

	case "file":
		ingest, err = store.NewLocalFile("./file/" + xmlbasename + ".json")
		if err != nil {
			panic(err)
		}
		defer ingest.(*store.LocalFile).FlushClose()

	default:
		fPln("[-store] is needed and from [map kvdb/badger file]")
		return
	}

	// -------------------- //

	xt.SetSlim(true)
	xt.SetIgnrAttr(
		"xsi:nil",
		"xmlns:xsd",
		"xmlns:xsi",
		"xmlns",
	)
	// xt.SetSuffix4List(`List`)    // xml-tool v0.1.18
	// xt.SetListPathSuffix(`List`) // xml-tool v0.1.23+
	xt.SetAttrPrefix(attrPrefix)
	xt.SetContAttrName(contAttrName)

	if err := xt.SetPathByFile("LIST", "./x2j-info/LIST.txt", "listSIF347", true, '/'); err != nil {
		panic(err)
	}
	if err := xt.EnableListPath("listSIF347"); err != nil {
		panic(err)
	}

	if err := xt.SetPathByFile("TYPE", "./x2j-info/BOOLEAN.txt", "typeSIF347", true, '/'); err != nil {
		panic(err)
	}
	if err := xt.SetPathByFile("TYPE", "./x2j-info/NUMERIC.txt", "typeSIF347", true, '/'); err != nil {
		panic(err)
	}
	if err := xt.EnableNonStrPath("typeSIF347"); err != nil {
		panic(err)
	}

	fPf("[%s] - [%d] elements saved into [%s]\n", *xmlPathPtr, scan(*xmlPathPtr, true, *asyncPtr, ingest), *storeTypePtr)
}
