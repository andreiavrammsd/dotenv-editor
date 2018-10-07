.PHONY: build test qa

BUILD := $(CURDIR)/build
BIN := $(BUILD)/dotenv

all: build

install:
	go get -t ./...
	go get -u \
	    github.com/go-bindata/go-bindata/... \
	    github.com/tdewolff/minify/cmd/minify
	curl -L https://git.io/vp6lP | sh
	sudo mv bin/* /usr/local/bin && rm -r bin
	sudo apt install upx -y

run:
	go-bindata -debug ui/...
	go run .

test:
	go test ./...

qa:
	gometalinter \
		--enable=staticcheck \
		--enable=gosimple \
		--enable=gochecknoglobals \
		--enable=gofmt \
		--enable=gochecknoinits \
		--enable=goimports \
		--enable=lll \
		--enable=nakedret \
		--enable=unparam \
		--enable=unused \
		./...

build: test
	mkdir -p $(BUILD)

	cp -r ui $(BUILD)
	minify -o $(BUILD)/ui ui/
	go-bindata -prefix="$(BUILD)/" $(BUILD)/ui/...

	@ # https://gist.github.com/FiloSottile/6c41098cc238988a900edb2ec5d27e6f
	GOOS=linux go build -ldflags="-s -w" -o $(BIN).tmp
	upx -f --brute -o $(BIN) $(BIN).tmp
	rm $(BIN).tmp

clean:
	rm -r build
