# ---------- Build ----------
FROM golang:1.26.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o smartcar-backend .

# ---------- Runtime ----------
FROM alpine:3.22

RUN apk add --no-cache wget

WORKDIR /app

COPY --from=builder /app/smartcar-backend .

EXPOSE 8080

CMD ["./smartcar-backend"]