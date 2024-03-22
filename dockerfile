## Build
FROM golang:1.22.1-alpine AS build

WORKDIR $GOPATH/src/segokuning

# manage dependencies
COPY . .
RUN go mod download

RUN go build -a -o /segokuning-server ./main.go


## Deploy
FROM alpine:latest
RUN apk add tzdata
COPY --from=build /segokuning-server /segokuning-server

EXPOSE 8080

ENTRYPOINT ["/segokuning-server"]