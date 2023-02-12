FROM golang:1.19-bullseye

# 1. Install the TensorFlow C Library (v2.11.0).
RUN curl -L https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-cpu-linux-$(uname -m)-2.11.0.tar.gz \
    | tar xz --directory /usr/local \
    && ldconfig

# 2. Install the Protocol Buffers Library and Compiler.
RUN apt-get update && apt-get -y install --no-install-recommends \
    libprotobuf-dev \
    protobuf-compiler

# 3. Install and Setup the TensorFlow Go API.
RUN git clone --branch=v2.11.0 https://github.com/tensorflow/tensorflow.git /go/src/github.com/tensorflow/tensorflow \
    && cd /go/src/github.com/tensorflow/tensorflow \
    && go mod init github.com/tensorflow/tensorflow \
    && sed -i '72 i \    ${TF_DIR}\/tensorflow\/tsl\/protobuf\/*.proto \\' tensorflow/go/genop/generate.sh \
    && (cd tensorflow/go/op && go generate) \
    && go mod edit -require github.com/google/tsl@v0.0.0+incompatible \
    && go mod edit -replace github.com/google/tsl=/go/src/github.com/google/tsl \
    && (cd /go/src/github.com/google/tsl && go mod init github.com/google/tsl && go mod tidy) \
    && go mod tidy \
    && go test ./...

# Build the Mentha Program.
WORKDIR /mentha
COPY main.go .
RUN go mod init app \
    && go mod edit -require github.com/google/tsl@v0.0.0+incompatible \
    && go mod edit -require github.com/tensorflow/tensorflow@v2.11.0+incompatible \
    && go mod edit -replace github.com/google/tsl=/go/src/github.com/google/tsl \
    && go mod edit -replace github.com/tensorflow/tensorflow=/go/src/github.com/tensorflow/tensorflow \
    && go mod tidy \
    && go build


ENTRYPOINT ["/mentha/app"]

