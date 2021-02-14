FROM golang:1.13 as builder

# modules
WORKDIR $GOPATH/src/github.com/shanehowearth/lb/racing
ADD . $GOPATH/src/github.com/shanehowearth/lb/racing

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# build time
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /go/bin/racing cmd/main.go

FROM scratch
WORKDIR /root/
COPY --from=0 /go/bin/racing .
# run options
ENTRYPOINT ["./racing"]
