FROM golang:1.20 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM alpine AS run

WORKDIR /

COPY --from=build /go/bin/app /app

ENTRYPOINT ["/app"]
