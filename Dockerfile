FROM golang:1.20 AS build

WORKDIR /app

COPY ./ ./
RUN go mod download

RUN CGO_ENABLED=0 go build -o /go/bin/app ./cmd/ynabnotifier/ynabnotifier.go

FROM alpine AS run

WORKDIR /

COPY --from=build /go/bin/app /app

ENTRYPOINT ["/app"]
