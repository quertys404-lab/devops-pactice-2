# Этап 1: сборка приложения в образе с Go.
FROM golang:1.22-alpine AS builder
WORKDIR /build

# Сначала скачиваем зависимости: этот слой можно переиспользовать в кэше.
COPY go.mod ./
RUN go mod download

COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/app ./main.go

# Этап 2: минимальный runtime без компилятора и shell.
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=builder --chown=nonroot:nonroot /out/app /app/app
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/app"]
