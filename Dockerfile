FROM alpine:3.16

# Dev
RUN apk update && \
    apk upgrade && \
    apk add --no-cache curl bash git jq yq && \
    apk add --no-cache pkgconfig gcc gpgme-dev libc-dev libcurl

# Python
RUN apk add --no-cache python3 py3-pip python3-dev && \
    ln -sf /usr/bin/python3 /usr/bin/python && \
    pip install awscli aws-sam-cli

# NodeJS
RUN apk add --no-cache nodejs npm && \
    npm config set unsafe-perm true && \
    npm update -g && \
    npm install -g aws-cdk

# Golang
COPY --from=golang:1.19-alpine /usr/local/go/ /usr/local/go/
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/go"

# Cevixe executables
COPY bin /cevixe/bin
ENV PATH="/cevixe/bin:${PATH}"

# Cevixe workspace

COPY sdk /cevixe/sdk
WORKDIR /cevixe/sdk
RUN go mod download

COPY app /cevixe/app
WORKDIR /cevixe/app
RUN go mod download

COPY cdk /cevixe/cdk

WORKDIR /cevixe/cdk
RUN go mod download

COPY test/app /cevixe/mod

# Cevixe workspace
ENV CEVIXE_SDK_HOME="/cevixe/sdk"
ENV CEVIXE_CDK_HOME="/cevixe/cdk"
ENV CEVIXE_APP_HOME="/cevixe/app"

WORKDIR /cevixe/mod