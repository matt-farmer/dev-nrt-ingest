# dev-nrt-ingest

An implementation for NRT ingestion.
Input a report XML file, convert each object to JSON and save them in kvdb(badger), or file.

## how to use

1. `go build` to create executable. Default executable name is 'dev-nrt-ingest'.

2. An usage example: `./dev-nrt-ingest -input=rrd.xml -check=false -store=kvdb -bar=false`.
   `rrd.xml` is input xml file path. One sample can be unzipped from rrdxml.zip. (you can also use sifxml.zip or your own report xml).
   `kvdb` is storage type, option is from [kvdb file map].
   `-check` switch is for debugging, which forces program to validate each object's XML and JSON.
   `-bar` switch is to show/hide processing progress.

3. BadgerDB is chosen for kvdb in this project. Storage files are under ./db if `-store=kvdb`

4. Storage JSON file is under ./file if `-store=file`

### one play report

HW: CPU: i3-9100f,  RAM: 16GB,  HD: 128GB SSD
OS: Ubuntu 20.04.1 LTS
XML: rrd.xml, No debugging check

no bar, store: map
time: 4.7s

with bar, store: map
time: 5.5s

no bar, store: kvdb
time: 5.4s

with bar, store: kvdb
time: 6.1s

no bar, store: file
time: 5.7s

with bar, store: file
time: 7.3s
