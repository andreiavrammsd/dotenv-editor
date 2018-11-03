.PHONY: build test qa

BUILD := $(CURDIR)/build
BIN := $(BUILD)/dotenv-editor
COVER := $(BUILD)/cover.out
IMAGE := andreiavrammsd/dotenv-editor

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

run: bindatadebug
	go run .

test: bindatadebug
	go test ./...

convey: init bindatadebug
	goconvey

qa: init
	gometalinter --config dev/.gometalinter.json ./...

build: init test qa
	cp -r ui $(BUILD)
	minify -o $(BUILD)/ui ui/
	go-bindata -pkg="handlers" -o="./handlers/bindata.go" -prefix="$(BUILD)/" $(BUILD)/ui/...

	@ # https://gist.github.com/FiloSottile/6c41098cc238988a900edb2ec5d27e6f
	GOOS=linux go build -ldflags="-s -w" -o $(BIN).tmp
	upx -f --brute -o $(BIN) $(BIN).tmp
	rm $(BIN).tmp

clean:
	@ if [ -d $(BUILD) ]; then rm -r $(BUILD); fi

init:
	@ mkdir -p $(BUILD)

bindatadebug:
	go-bindata -debug -pkg="handlers" -o="./handlers/bindata.go" ui/...

dockerbuild:
	docker build . -f dev/Dockerfile -t $(IMAGE)

dockerpush:
	docker push $(IMAGE)

dockertestqa:
	docker pull $(IMAGE)
	docker run -ti --rm -v $(CURDIR):/go/src/github.com/andreiavrammsd/dotenv-editor $(IMAGE) make test qa
