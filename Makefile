GOBIN = ./build/bin

run-concurrent-set: clean build-concurrency-set
	$(GOBIN)/concurrent-set
build-concurrency-set:
	go build -o $(GOBIN)/concurrent-set ./cmd/go-concurrency/ConcurrentSet

run-concurrent-map: clean build-concurrency-map
	$(GOBIN)/concurrent-map
build-concurrency-map:
	go build -o $(GOBIN)/concurrent-map ./cmd/go-concurrency/ConcurrentMap

clean:
	rm -fr $(GOBIN)/*

revive:
	$(GOPATH)/bin/revive -config ./revive.toml -formatter friendly ./...