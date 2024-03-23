FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o main ./cmd/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/main /app
COPY --from=builder /app/.env /app
COPY --from=builder /app/firebase-service-account-key.json /app/firebase-service-account-key.json
CMD [ "./main" ]