FROM golang:1.13 AS build
ENV GO111MODULE=on
WORKDIR /github.com/olivebay/urlinfo
COPY go.mod ./
COPY . /github.com/olivebay/urlinfo/
RUN go mod download
RUN CGO_ENABLED=0 go build -o urlinfo-api

FROM scratch
WORKDIR /app
COPY --from=build /github.com/olivebay/urlinfo/urlinfo-api /app/

EXPOSE 9090
CMD ["/app/urlinfo-api"]
