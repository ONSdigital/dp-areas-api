module github.com/ONSdigital/dp-topic-api

go 1.15

require (
	github.com/ONSdigital/dp-api-clients-go v1.32.5
	github.com/ONSdigital/dp-authorisation v0.1.0
	github.com/ONSdigital/dp-component-test v0.2.0
	github.com/ONSdigital/dp-healthcheck v1.0.5
	github.com/ONSdigital/dp-mongodb v1.5.0
	github.com/ONSdigital/dp-net v1.0.11
	github.com/ONSdigital/go-ns v0.0.0-20200902154605-290c8b5ba5eb // indirect
	github.com/ONSdigital/log.go v1.0.1
	github.com/benweissmann/memongo v0.1.1
	github.com/cucumber/godog v0.10.0
	github.com/fatih/color v1.10.0 // indirect; indirectfind.go:115:1
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/hokaccha/go-prettyjson v0.0.0-20210113012101-fb4e108d2519 // indirect
	github.com/justinas/alice v1.2.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/smartystreets/goconvey v1.6.4
	go.mongodb.org/mongo-driver v1.5.2
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777 // indirect
	golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22 // indirect
)

replace github.com/ONSdigital/dp-mongodb v1.5.0 => github.com/ONSdigital/dp-mongodb v1.5.1-0.20210613162655-0f90bde40748
