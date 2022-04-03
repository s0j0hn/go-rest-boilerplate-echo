# Base image:
FROM golang:1.18-alpine3.15

# Install clang from LLVM repository
RUN apk update && apk add --no-cache clang alpine-sdk make
