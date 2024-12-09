FROM golang:1.23 AS builder

# 设置工作目录
WORKDIR /app

COPY ./ /app/

RUN apt-get update && \
apt-get install -y unzip && \
touch ./internal/alyzers/router/static/alyzers.js && \
make -f build/Makefile build

FROM golang:1.23

RUN apt-get install -y tzdata && \
mkdir -p /opt/alyzers/bin /opt/alyzers/conf.d

COPY --from=builder /app/alyzers /opt/alyzers/bin

EXPOSE 8080

WORKDIR /opt/alyzers

ENTRYPOINT [ "./bin/alyzers -conf conf.d/config.toml" ]
