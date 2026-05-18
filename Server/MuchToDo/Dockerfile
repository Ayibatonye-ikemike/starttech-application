FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates curl
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
ENV PORT=8080
CMD ["./main"]
# manual-kick-1779112049
# fix deployment loop
# trigger execution buffer test
# authenticated production database release
# tracking production version rollout
# fix container boot loop sequence
# diagnostic validation patch
# final production rollout verification
