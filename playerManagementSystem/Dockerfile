# build stage
FROM golang:1.20 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN GOSUMDB=off go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOPROXY=direct GOSUMDB=off go build -o server .

# final stage
FROM alpine:3.12 as production

COPY --from=builder /app/server /app/server

CMD ["/app/server"]