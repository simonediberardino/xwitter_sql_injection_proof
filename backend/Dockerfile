FROM golang:1.23-alpine AS build

WORKDIR /app
COPY go.mod ./
COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o /xwitter-api ./cmd/server

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=build /xwitter-api /xwitter-api
EXPOSE 4000
USER nobody
CMD ["/xwitter-api"]
