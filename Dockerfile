FROM golang:1.12.7-alpine as builder

RUN apk update
RUN apk add git gcc

WORKDIR /go/src/github.com/form-is-function/roscetta
COPY Gopkg.* ./

RUN go get -u github.com/golang/dep/cmd/dep

RUN dep ensure -v -vendor-only

COPY . .

ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7
ENV CGO_ENABLED=0
RUN go build -o /tmp/proxy cmd/proxy/main.go

FROM balenalib/armv7hf-alpine

WORKDIR /opt/roscetta

COPY --from=builder /tmp/proxy /opt/roscetta/

CMD ["/opt/roscetta/proxy"]
