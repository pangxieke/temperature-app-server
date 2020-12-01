FROM registry.cn-shenzhen.aliyuncs.com/yunzhimeng/portrays-builder:latest as builder
WORKDIR /go/src/temperature

COPY . .
ENV GOPROXY https://goproxy.io
ENV GO111MODULE on
RUN go mod download

#RUN go test ./... -coverprofile .testCoverage.txt \
#    && go tool cover -func=.testCoverage.txt
RUN CGO_ENABLED=0 go build -o app_d ./cmd/main.go \
     && CGO_ENABLED=0 go build ./cmd/migrate

FROM alpine:3.8
RUN apk --no-cache add ca-certificates
LABEL \
    SERVICE_80_NAME=temperature_http \
    SERVICE_NAME=temperature \
    description="temperature" \
    maintainer="汤良浩"

EXPOSE 8080
COPY --from=builder /go/src/temperature/app_d /bin/app
COPY --from=builder /go/src/temperature/migrate /bin/migrate
ENTRYPOINT ["app"]
