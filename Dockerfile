# Base image:
FROM golang:1.14-alpine

# Install clang from LLVM repository
RUN apk update && apk add --no-cache clang make gcc libc-dev
# Set Clang as default CC
ENV set_clang /etc/profile.d/set-clang-cc.sh
RUN echo "export CC=clang-9" | tee -a ${set_clang} && chmod a+x ${set_clang}
RUN uname -r; echo $(gcc -v 2>&1 | grep version); echo $(clang -v 2>&1 | grep version); go version
