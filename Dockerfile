# Base image:
FROM golang:1.15-alpine

# Install clang from LLVM repository
RUN apk update && apk add --no-cache clang make gcc g++ libc-dev
