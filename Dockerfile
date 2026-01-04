ARG GO_VERSION=1.25.5
FROM golang:${GO_VERSION} as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux \
    CGO_ENABLED=0 \
    go build \
    -ldflags="-s -w" \
    -o server ./cmd/server/main.go

FROM alpine:3.23.2

# Instala ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates

# Cria usuário não-privilegiado
RUN addgroup -g 1000 appgroup && \
    adduser -D -u 1000 -G appgroup appuser

WORKDIR /app

# Copia o binário com permissões corretas
COPY --from=builder --chown=appuser:appgroup /build/server .

# Muda para o usuário não-privilegiado
USER appuser

# Expõe a porta
EXPOSE 8080

# Modo padrão: api
# No Kubernetes, override com: args: ["kafka"] ou ["rabbitmq"] ou ["grpc"]
CMD [ "./server", "api" ]
