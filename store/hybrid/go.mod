module github.com/kshard/spock/store/hybrid

go 1.20

require (
	github.com/fogfish/faults v0.2.0
	github.com/fogfish/segment v0.0.0-20230508170851-362615e751a8
	github.com/fogfish/skiplist v0.12.0
	github.com/kshard/spock v0.2.0
	github.com/kshard/xsd v0.1.0
)

require (
	github.com/fogfish/curie v1.8.2 // indirect
	github.com/fogfish/guid/v2 v2.0.4 // indirect
)

replace github.com/fogfish/segment => ../../../../fogfish/segment

replace github.com/fogfish/skiplist => ../../../../fogfish/skiplist

replace github.com/fogfish/golem/trait => ../../../../fogfish/golem/trait
