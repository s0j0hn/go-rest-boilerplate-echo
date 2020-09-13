# Base image:
FROM golang:1.15-alpine

# Install clang from LLVM repository
RUN apk add alpine-sdk curl
RUN apk update && apk add \
  && rm -rf /var/lib/apt/lists/* \
  && curl -SL https://github.com/llvm/llvm-project/releases/download/llvmorg-10.0.0/clang+llvm-10.0.0-x86_64-linux-gnu-ubuntu-18.04.tar.xz \
  | tar -xJC . && \
  mv clang+llvm-10.0.0-x86_64-linux-gnu-ubuntu-18.04 clang_10.0.0 && \
  echo 'export PATH=/clang_10.0.0/bin:$PATH' >> ~/.bashrc && \
  echo 'export LD_LIBRARY_PATH=/clang_10.0.0/lib:$LD_LIBRARY_PATH' >> ~/.bashrc

# Start from a Bash prompt
CMD [ "/bin/bash" ]
