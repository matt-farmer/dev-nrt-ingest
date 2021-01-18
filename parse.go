package main

import (
	"io/ioutil"
	"os"

	jt "github.com/cdutwhu/json-tool"
	xt "github.com/cdutwhu/xml-tool"
)

const (
	sep, preatt = "~", "-"
)

func readEleLn(xml string) string {
	s := -1
	for i, c := range xml {
		if s == -1 && c != '\n' && c != '\t' {
			s = i
			continue
		}
		if s != -1 && c == '\n' && xml[i-1] == '>' {
			return xml[s:i]
		}
	}
	return ""
}

// ONLY for formatted one line
func lineTCAV(xml string) (tag, cont string, attrs []string, mAttrVal map[string]string) {
	s, e := -1, -1
	hasAttr := false

	// tag
	for i, c := range xml {
		if c == '<' {
			s = i + 1
		}
		if c == ' ' || c == '>' || c == '/' {
			e = i
			if c == ' ' {
				hasAttr = true
			}
			break
		}
	}
	if e == -1 {
		return
	}

	tag = xml[s:e]
	xml = xml[e+1:]

	// attributes & content
	if hasAttr {
		mAttrVal = make(map[string]string)
		eql := -1
		attr := ""
	AGAIN:
		sv, ev := -1, -1
		for i, c := range xml {
			if c == '=' {
				eql = i
				attr = sTrimLeft(xml[:eql], " ")
				continue
			}
			if c == '"' {
				if i-eql == 1 {
					sv = i + 1
				} else {
					ev = i
				}
			}
			if sv != -1 && ev != -1 {
				attrs = append(attrs, attr)
				mAttrVal[attr] = xml[sv:ev]
				xml = xml[ev+1:] // next part, from ' ' or '>' or '/'
				goto AGAIN
			}
		}
	}

	// content :
	if hasAttr {
		xml = xml[1:] // skip '>' OR ' ' of " />"
	}
	if len(xml) > 0 && xml[0] == '/' { // empty content
		return
	}
	for i, c := range xml {
		if c == '<' {
			cont = sTrimRight(xml[:i], " \t\n")
			break
		}
	}

	return
}

func isCplxSLn(ln string) bool {
	return !sHasSuffix(ln, "/>") && !sContains(ln, "</")
}

func parse(xml string, ingest IIngest) {

	// FIRST (root) line
	ln := readEleLn(xml)
	tag, cont, as, mav := lineTCAV(ln)

	RefID, ok := mav[attrNameOfID]
	if !ok {
		panic("error sif object")
	}

	ingest.AddVal(RefID, []byte(xml), false) // SAVE a copy of object whole string ... ZIP (store/map.go)

	path := RefID + "@" + tag      // uuid@root
	ipath := ingest.GenIPath(path) // ipath := getIPath(path)

	basePath := make(map[int]string)
	basePath[0] = ipath // update basePath, key:lvl, value:ipath

	// content (complex)
	if isCplxSLn(ln) {
		// dbgPln(false, "CPLX :", ipath, xml)
	} else { // content (simple)
		ingest.AddVal(ipath, []byte(cont), false)
		// dbgPln(false, "CONT :", ipath, cont)
	}

	// attributes
	for _, a := range as {
		ingest.AddVal(ipath+sep+preatt+a, []byte(mav[a]), false)
		// dbgPln(false, "ATTR :", ipath+sep+preatt+a, mav[a])
	}
}

func parse4json(xml string, ingest IIngest) {

	errRecord := func(xml, jstr string) {
		ioutil.WriteFile("err-json.xml", []byte(xml), os.ModePerm)
		ioutil.WriteFile("err-json.json", []byte(jstr), os.ModePerm)
	}

	// DEBUG checking ...
	if validate {
		if !xt.IsValid(xml) {
			errRecord(xml, "")
			panic("*ERROR*, Invalid XML")
		}
	}

	// xml -> json
	jstr := xt.MkJSON(xml)
	RefID, ok := jt.SglEleAttrVal(jstr, attrNameOfID, attrPrefix)
	if !ok {
		errRecord(xml, jstr)
		panic("ERROR, NO RefId found")
	}

	// DEBUG checking ...
	if validate {
		if !jt.IsValid(jstr) {
			errRecord(xml, jstr)
			panic("*ERROR*, Invalid JSON")
		}
	}

	ingest.AddVal(RefID, []byte(jstr), false) // SAVE a copy of object string ... ZIP (store/map.go)
}
