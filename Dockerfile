FROM golang:latest as builder
#RUN apk add --no-cache curl 
#RUN apk add --no-cache git 
RUN mkdir -p /go/src/github.com/jiseruk/minesweeper
WORKDIR /go/src/github.com/jiseruk/minesweeper
COPY go.mod go.sum ./
#RUN go list -e $(go list -f '{{.Path}}' -m all 2>/dev/null)
RUN go mod download
RUN GO111MODULE=on go mod vendor
#For local environment
#ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh ./
#RUN chmod +x ./wait-for-it.sh
RUN go get -u github.com/swaggo/swag/cmd/swag
COPY . .
WORKDIR /go/src/github.com/jiseruk/minesweeper/cmd/minesweeper
RUN swag init

FROM builder as tests
ENV GOPATH /go
WORKDIR /go/src/github.com/jiseruk/minesweeper
#RUN go test ./... -covermode=count -coverprofile=cover.out -coverpkg=./...
WORKDIR /go/src/github.com/jiseruk/minesweeper/cmd/minesweeper
#RUN bin/tests.sh
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main . 

FROM alpine:latest
RUN apk update && apk add bash
ENV GOPATH /go
ENV GIN_MODE release
WORKDIR /go/src/github.com/jiseruk/minesweeper/cmd/minesweeper
COPY --from=tests /go/src/github.com/jiseruk/minesweeper/wait-for-it.sh /go/src/github.com/jiseruk/minesweeper/cmd/minesweeper/wait-for-it.sh
COPY --from=tests /go/src/github.com/jiseruk/minesweeper/cmd/minesweeper/main /go/src/github.com/jiseruk/minesweeper/cmd/minesweeper/main 
COPY --from=tests /go/src/github.com/jiseruk/minesweeper/cmd/minesweeper/config/local.yaml /go/src/github.com/jiseruk/minesweeper/cmd/minesweeper/config/local.yaml
COPY --from=tests /go/src/github.com/jiseruk/minesweeper/cmd/minesweeper/config/now.yaml /go/src/github.com/jiseruk/minesweeper/cmd/minesweeper/config/now.yaml
ENTRYPOINT ["./main"]
#CMD ["./main"]
EXPOSE 8080 
#CMD ["go", "run", "main.go"]
