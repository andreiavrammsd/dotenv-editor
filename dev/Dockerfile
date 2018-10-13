FROM golang:1-alpine

RUN apk add --update --no-cache alpine-sdk make git curl

COPY Makefile .

RUN go get -u \
	    github.com/go-bindata/go-bindata/... \
	    github.com/smartystreets/goconvey \
	    github.com/stretchr/testify/mock \
	    github.com/tdewolff/minify/cmd/minify \
	&& curl -L https://git.io/vp6lP | sh \
	&& mv bin/* /usr/local/bin && rm -r bin

WORKDIR /go/src/github.com/andreiavrammsd/dotenv-editor