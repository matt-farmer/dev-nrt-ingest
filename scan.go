package main

import (
	"bufio"
	"os"
	"sync"

	xt "github.com/cdutwhu/xml-tool"
	"github.com/gosuri/uiprogress"
)

// scan :
func scan(xmlpath string, cvt2json, async bool, ingest IIngest) (count uint64) {

	file, err := os.Open(xmlpath)
	if err != nil {
		panic(err)
	}
	fileinfo, err := file.Stat()
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// -------------------------- //

	// progress bar
	if probar {
		uip = uiprogress.New()
		defer uip.Stop()
		uip.Start()
		bar = uip.AddBar(int(fileinfo.Size()))
		bar.AppendCompleted().PrependElapsed()
	}

	// -------------------------- //

	br := bufio.NewReader(file)
	var dataTypes = []string{
		"NAPStudentResponseSet",
		"NAPEventStudentLink",
		"StudentPersonal",
		"NAPTestlet",
		"NAPTestItem",
		"NAPTest",
		"NAPCodeFrame",
		"SchoolInfo",
		"NAPTestScoreSummary",
	}

	if async {
		var wg sync.WaitGroup
		for ele := range xt.StreamEle(br, dataTypes...) {
			count++
			go proc(&wg, count, ele, cvt2json, ingest)
		}
		wg.Wait()
	} else {
		for ele := range xt.StreamEle(br) {
			count++
			proc(nil, count, ele, cvt2json, ingest)
		}
	}

	return
}
