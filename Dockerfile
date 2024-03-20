FROM golang:alpine AS build-env
LABEL MAINTAINER "vivek singh"
ENV GOPATH /go
WORKDIR /go/src
COPY . /go/src/scraper
RUN cd /go/src/scraper && go build -o scraper main.go

FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk*
WORKDIR /app
COPY --from=build-env /go/src/scraper /app
COPY .env /app

EXPOSE 8080

ENTRYPOINT [ "./scraper" ]