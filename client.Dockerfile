FROM golang:1.24 AS build

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go mod download
RUN go test ./...

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o client ./cmd/client/client.go

FROM gcr.io/distroless/static-debian11

COPY --from=build /app/client /client

COPY config/client/config.yaml /config/config.yaml

ENTRYPOINT ["/client"]
