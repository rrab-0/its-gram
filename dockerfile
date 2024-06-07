FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o main ./cmd/main.go
# RUN go build -o main ./cmd/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/main /app
COPY --from=builder /app/.env /app
COPY --from=builder /app/firebase-service-account-key.json /app/firebase-service-account-key.json
COPY --from=builder /app/web /app/web
ENTRYPOINT [ "./main", "DEV" ]