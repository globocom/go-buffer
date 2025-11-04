module github.com/globocom/go-buffer/v4

retract (
	v3.0.0 // Published prematurely
	v3.0.1 // Contains retractions only
)

go 1.24

retract (
	v3
)

require (
	github.com/onsi/ginkgo v1.13.0
	github.com/onsi/gomega v1.10.1
)

require (
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/nxadm/tail v1.4.4 // indirect
	golang.org/x/net v0.0.0-20200520004742-59133d7f0dd7 // indirect
	golang.org/x/sys v0.0.0-20200519105757-fe76b779f299 // indirect
	golang.org/x/text v0.3.2 // indirect
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)
