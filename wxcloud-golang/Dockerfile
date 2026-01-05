# 选择构建用基础镜像
FROM golang:1.24-alpine as builder

# 指定构建过程中的工作目录
WORKDIR /app

# 【关键修改 1】设置国内代理，解决依赖包下载超时/失败的问题
ENV GOPROXY=https://goproxy.cn,direct
# 开启 Go Modules 模式
ENV GO111MODULE=on

# 【关键修改 2】先只拷贝 go.mod 和 go.sum，利用 Docker 缓存
COPY go.mod go.sum ./

# 下载依赖（如果 go.mod 没变，这一步会直接使用缓存，极快）
RUN go mod download

# 将当前目录下所有文件拷贝到工作目录下
COPY . .

# 执行代码编译命令
# CGO_ENABLED=0 确保编译出纯静态二进制文件，避免依赖库问题
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 选用运行时所用基础镜像
FROM alpine:3.13

# 安装证书（访问 HTTPS 需要）和时区数据
RUN apk add --no-cache ca-certificates tzdata

# 【可选】设置时区为上海（推荐打开，否则日志时间会少8小时）
ENV TZ=Asia/Shanghai

# 指定运行时的工作目录
WORKDIR /app

# 【关键修改 3】拷贝二进制文件
# 注意：我去掉了 index.html，因为你做 API 服务可能不需要它。
# 如果你确实有 index.html 文件，请把下面这行改成: COPY --from=builder /app/main /app/index.html /app/
COPY --from=builder /app/main /app/

# 执行启动命令
CMD ["/app/main"]