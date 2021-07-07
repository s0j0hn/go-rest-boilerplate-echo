# Base image:
FROM golang:1.16-alpine

# Install clang from LLVM repository
RUN apk update && apk add --no-cache clang alpine-sdk make
