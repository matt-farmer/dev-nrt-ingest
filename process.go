package main

import (
	"sync"
	"sync/atomic"

	xt "github.com/cdutwhu/xml-tool"
)

var mutex = &sync.Mutex{}

func proc(params ...interface{}) error {

	var (
		wg       *sync.WaitGroup
		id       = params[1].(uint64)  // threadID
		xml      = params[2].(string)  // input xml
		cvt2json = params[3].(bool)    // whether to convert to json
		ingest   = params[4].(IIngest) // ingest interface
	)

	if params[0] != nil {
		wg = params[0].(*sync.WaitGroup) // WaitGroup
		defer wg.Done()
		wg.Add(1)
	}

	if probar {
		// mutex.Lock()
		// -- progress bar -- //
		atomic.AddUint64(&procsize, uint64(len(xml)))
		bar.Set(int(procsize))
		bar.Incr()
		// mutex.Unlock()
	}

	dbgPln(false, "---@---", id)

	xml = xt.RmEmptyEle(xml, 3, false)
	if cvt2json {
		parse4json(xml, ingest)
	} else {
		parse(xt.Fmt(xml), ingest) // if store xml, Fmt it. if store json, do not need fmt
	}
	return nil
}
