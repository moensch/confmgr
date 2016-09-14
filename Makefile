build:
	go build -o ./bin/confmgr ./cmd/confmgr

test:
	go test -v ./tests

coverage:
	go get github.com/jstemmer/go-junit-report
	go get github.com/axw/gocov/gocov
	go get github.com/AlekSi/gocov-xml
	go test -v -covermode=count -coverprofile=coverage.out ./tests | go-junit-report > test_results.xml
	gocov convert coverage.out | gocov-xml > coverage_cobertura.xml
	go tool cover -html=coverage.out -o coverage.html

run: build
	./bin/confmgr
