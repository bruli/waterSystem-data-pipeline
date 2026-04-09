# syntax=docker/dockerfile:1.20

############################
# Etapa de build ARM64
############################
FROM --platform=$BUILDPLATFORM golang:1.26.2 AS builder
WORKDIR /src

ENV GOPROXY=https://proxy.golang.org,direct

# 1) Deps (capa estable + cache)
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# 2) Codi
COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/waterSystemDataPipeline ./cmd/waterSystemDataPipeline
# --- runtime ---
FROM alpine:3.22

RUN apk add --no-cache tzdata
# copia el teu binari
COPY --from=builder /out/waterSystemDataPipeline /usr/local/bin/waterSystemDataPipeline

ENTRYPOINT ["/usr/local/bin/waterSystemDataPipeline"]
