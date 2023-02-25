FROM golang:1.19 as builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
ENV GOPROXY=https://goproxy.cn
RUN go mod download

COPY main.go main.go
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -ldflags "-s -w -extldflags -static" -o request main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/request .
USER 65532:65532

ENTRYPOINT ["/request"]
