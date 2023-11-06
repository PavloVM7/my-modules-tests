GOBIN = ./build/bin

run-go-concurrency: clean build-go-concurrency
	$(GOBIN)/go-concurrency
build-go-concurrency:
	go build -o $(GOBIN)/go-concurrency ./go-concurrency

clean:
	rm -fr $(GOBIN)/*

revive:
	$(GOPATH)/bin/revive -config ./revive.toml -formatter friendly ./...