# Base image:
FROM golang:1.17-alpine3.14

# Install clang from LLVM repository
RUN apk update && apk add --no-cache clang alpine-sdk make
