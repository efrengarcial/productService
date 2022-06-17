.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/products products/cmd/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
