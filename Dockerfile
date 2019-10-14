FROM golang:latest as builder
#RUN apk add --no-cache curl 
#RUN apk add --no-cache git 
RUN mkdir -p /go/src/github.com/jiseruk/minesweeper
WORKDIR /go/src/github.com/jiseruk/minesweeper
COPY go.mod go.sum ./
RUN GO111MODULE=on go mod vendor
#For local environment
ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh ./
RUN chmod +x ./wait-for-it.sh
RUN go get -u github.com/swaggo/swag/cmd/swag
COPY . .

FROM builder as tests
ENV GOPATH /go
WORKDIR /go/src/github.com/jiseruk/minesweeper
#RUN go test ./... -covermode=count -coverprofile=cover.out -coverpkg=./...
WORKDIR /go/src/github.com/jiseruk/minesweeper/cmd/minesweeper
RUN swag init
#RUN bin/tests.sh
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main . 

#FROM busybox:musl
FROM alpine:latest
ENV GOPATH /go
ENV GIN_MODE release
WORKDIR /go/src/github.com/jiseruk/minesweeper/cmd/minesweeper
COPY --from=tests /go/src/github.com/jiseruk/minesweeper/cmd/minesweeper/main /go/src/github.com/jiseruk/minesweeper/cmd/minesweeper/main 
COPY --from=tests /go/src/github.com/jiseruk/minesweeper/config/config.yaml /go/src/github.com/jiseruk/minesweeper/config/config.yaml
ENTRYPOINT ["./main"]
#CMD ["./main"]
EXPOSE 8080 
#CMD ["go", "run", "main.go"]
