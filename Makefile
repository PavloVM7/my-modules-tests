GOBIN = ./build/bin

run-go-concurrency-set: clean build-go-concurrency-set
	$(GOBIN)/go-concurrent-set
build-go-concurrency-set:
	go build -o $(GOBIN)/go-concurrent-set ./go-concurrency/ConcurrentSet

run-go-concurrency-map: clean build-go-concurrency-map
	$(GOBIN)/go-concurrent-map
build-go-concurrency-map:
	go build -o $(GOBIN)/go-concurrent-map ./go-concurrency/ConcurrentMap

clean:
	rm -fr $(GOBIN)/*

revive:
	$(GOPATH)/bin/revive -config ./revive.toml -formatter friendly ./...