# Build the manager binary
FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.15 as builder

RUN apt update && apt install unzip -y 

# ARG GOARCH=amd64
ENV CGO_ENABLED=0 
ENV GOOS=linux 
ENV GO111MODULE=on
ENV ONEPASSWORD_CLI_VERSION=v1.8.0

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

# install 1password binary
RUN cd /tmp; curl https://cache.agilebits.com/dist/1P/op/pkg/${ONEPASSWORD_CLI_VERSION}/op_linux_amd64_${ONEPASSWORD_CLI_VERSION}.zip -o op_linux_amd64_${ONEPASSWORD_CLI_VERSION}.zip; unzip op_linux_amd64_${ONEPASSWORD_CLI_VERSION}.zip; mv ./op /usr/local/bin/
RUN gpg --keyserver hkps://keyserver.ubuntu.com --recv-keys 3FEF9748469ADBE15DA7CA80AC2D62742012EA22
RUN cd /tmp; gpg --verify /tmp/op.sig /usr/local/bin/op || (echo "ERROR: Incorrect GPG signature for 1password op binary." && exit 1)


# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM --platform=${TARGETPLATFORM:-linux/amd64}  gcr.io/distroless/base-debian10
WORKDIR /
COPY --from=builder /workspace/manager .

COPY --from=builder /usr/local/bin/op  /usr/local/bin/ 

USER nonroot:nonroot

ENTRYPOINT ["/manager"]
