install:
	go install

release:
	goreleaser

setup:
	go get github.com/goreleaser/goreleaser
