# build stage
FROM golang:1.20 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN GOSUMDB=off go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOPROXY=direct GOSUMDB=off go build -o endlessChallengeSystem .

# final stage
FROM alpine:3.12 as production

COPY --from=builder /app/endlessChallengeSystem /app/endlessChallengeSystem

CMD ["/app/endlessChallengeSystem"]

# expose port 8085 for REST API server
EXPOSE 8085