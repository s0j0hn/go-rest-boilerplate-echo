# Base image:
FROM golang:1.19-alpine3.16

# Install clang from LLVM repository
RUN apk update && apk add --no-cache clang alpine-sdk make

RUN go install gotest.tools/gotestsum@latest