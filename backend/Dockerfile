# growth-partner/backend/Dockerfile
# 多阶段构建：development（热重载）+ production（精简镜像）
# 兼容 Mac (arm64/amd64) + Windows Docker Desktop (amd64)

# ─────────────────────────────────────────────────────────────
# 阶段一：Go 模块下载层（利用缓存，加速重建）
# ─────────────────────────────────────────────────────────────
FROM golang:1.23-alpine AS base

# 设置必要的系统依赖
RUN apk add --no-cache \
    git \
    curl \
    ca-certificates \
    tzdata \
    gcc \
    musl-dev

# 设置时区
ENV TZ=Asia/Shanghai
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

# 设置 Go 环境
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
ENV GOPATH=/go
ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /app

# 优先复制 go.mod/go.sum，利用 Docker 层缓存
# 只有依赖变化时才重新下载，代码变化不触发重下载
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# ─────────────────────────────────────────────────────────────
# 阶段二：开发环境（air 热重载）
# ─────────────────────────────────────────────────────────────
FROM base AS development

# 安装 air 热重载工具（版本固定，保证 Mac/Windows 行为一致）
RUN go install github.com/air-verse/air@v1.61.5

# air 配置文件
COPY .air.toml ./

# 开发模式不复制源码，通过 volume 挂载实现热重载
EXPOSE 8080 2345

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

CMD ["air", "-c", ".air.toml"]

# ─────────────────────────────────────────────────────────────
# 阶段三：生产构建
# ─────────────────────────────────────────────────────────────
FROM base AS builder

COPY . .

RUN go build \
    -ldflags="-w -s -X main.Version=$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')" \
    -o /build/growth-partner \
    ./main.go

# ─────────────────────────────────────────────────────────────
# 阶段四：生产运行镜像（最小化）
# ─────────────────────────────────────────────────────────────
FROM alpine:3.20 AS production

RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    curl

ENV TZ=Asia/Shanghai
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

# 创建非 root 用户运行（安全最佳实践）
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /build/growth-partner .
COPY --from=builder /app/migrations ./migrations

RUN chown -R appuser:appgroup /app
USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

CMD ["./growth-partner"]