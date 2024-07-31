FROM golang:1.22.3-alpine as builder

RUN apk --no-cache add make git

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...

RUN go build -ldflags="-w -s" -o /go/bin/result

FROM alpine

COPY --from=builder /go/bin/result /go/bin/result

ENV SKIP_DOWNLOAD=true
ENV VENDOR_PATH=/usr/bin/

EXPOSE 8080

CMD ["/go/bin/result"]
