# This script was copied from 'terraform-provider-akamai'
FROM golang:1.23.6-alpine3.21
ENV PROVIDER_VERSION="1.0.0" \
    GO111MODULE="on" \
    CGO_ENABLED=0 \
    GOOS="linux" \
    GOARCH="amd64" \
    PATH=$PATH:/root/go/bin

ARG SSH_PRV_KEY
ARG SSH_PUB_KEY
ARG SSH_KNOWN_HOSTS
ARG TERRAFORM_VERSION="1.9.5"
WORKDIR $GOPATH/src/github.com/akamai

RUN apk add --update git bash sudo openssh gcc musl-dev openssl-dev ca-certificates unzip curl make && \
    wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    mv terraform /usr/bin/terraform && \
    mkdir -p /root/.ssh

ADD build/AkamaiCorpRoot-G1.pem /usr/local/share/ca-certificates/AkamaiCorpRoot-G1.pem
RUN update-ca-certificates
