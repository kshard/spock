module github.com/kshard/spock/store/ephemeral

go 1.20

require (
	github.com/fogfish/curie v1.8.2
	github.com/fogfish/golem/trait v0.1.0
	github.com/fogfish/guid/v2 v2.0.2
	github.com/fogfish/it/v2 v2.0.1
	github.com/fogfish/skiplist v0.13.1
	github.com/kshard/spock v0.1.0
	github.com/kshard/xsd v0.1.0
)

replace github.com/kshard/xsd => ../../../../kshard/xsd

replace github.com/kshard/spock => ../../

replace github.com/fogfish/golem/trait => ../../../../fogfish/golem/trait
