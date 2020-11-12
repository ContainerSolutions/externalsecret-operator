# Build the manager binary
FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.15 as builder

RUN apt update && apt install unzip -y 

# ARG GOARCH=amd64
ENV CGO_ENABLED=0 
ENV GOOS=linux 
ENV GO111MODULE=on

ARG TARGETPLATFORM
RUN go env

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY apis/ apis/
COPY controllers/ controllers/
COPY pkg/ pkg/

# Build
RUN go build -a -o manager main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM --platform=${TARGETPLATFORM:-linux/amd64}  gcr.io/distroless/base-debian10@sha256:abe4b6cd34fed3ade2e89ed1f2ce75ddab023ea0d583206cfa4f960b74572c67
WORKDIR /
COPY --from=builder /workspace/manager .

USER nonroot:nonroot

ENTRYPOINT ["/manager"]
