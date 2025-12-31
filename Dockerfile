# Base image:
FROM golang:1.24-alpine3.21

# Install clang from LLVM repository
RUN apk update && apk add --no-cache clang alpine-sdk make

RUN go install gotest.tools/gotestsum@latest