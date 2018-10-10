.PHONY: build test qa

BUILD := $(CURDIR)/build
BIN := $(BUILD)/dotenv
COVER := $(BUILD)/cover.out

all: build

install:
	go get -u \
	    github.com/go-bindata/go-bindata/... \
	    github.com/smartystreets/goconvey \
	    github.com/stretchr/testify/mock \
	    github.com/tdewolff/minify/cmd/minify
	curl -L https://git.io/vp6lP | sh #gometalinter
	sudo mv bin/* /usr/local/bin && rm -r bin
	sudo apt install upx -y

run:
	go-bindata -debug -pkg="handlers" -o="./handlers/bindata.go" ui/...
	go run .

test:
	goconvey

qa:
	gometalinter \
		--enable=megacheck \
		--enable=gochecknoglobals \
		--enable=gofmt \
		--enable=gochecknoinits \
		--enable=goimports \
		--enable=lll \
		--enable=nakedret \
		--enable=unparam \
		./...

build: test
	mkdir -p $(BUILD)

	cp -r ui $(BUILD)
	minify -o $(BUILD)/ui ui/
	go-bindata -pkg="handlers" -o="./handlers/bindata.go" -prefix="$(BUILD)/" $(BUILD)/ui/...

	@ # https://gist.github.com/FiloSottile/6c41098cc238988a900edb2ec5d27e6f
	GOOS=linux go build -ldflags="-s -w" -o $(BIN).tmp
	upx -f --brute -o $(BIN) $(BIN).tmp
	rm $(BIN).tmp

clean:
	rm -r build
