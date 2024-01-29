FROM ubuntu:22.04

RUN apt-get update && apt-get -y upgrade
RUN apt-get -y install curl git
RUN apt-get -y install libpq-dev jq

RUN echo "alias ll='ls -ltr'" > /root/.bashrc

# Setup Rust
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y

# Setup Golang
RUN curl --location -o go.tgz https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go.tgz
ENV PATH $PATH:/usr/local/go/bin
RUN go version

RUN go install golang.org/x/tools/gopls@latest
RUN go install github.com/cweill/gotests/gotests@v1.6.0
RUN go install github.com/fatih/gomodifytags@v1.16.0
RUN go install github.com/josharian/impl@v1.1.0
RUN go install github.com/haya14busa/goplay/cmd/goplay@v1.0.0
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install honnef.co/go/tools/cmd/staticcheck@latest

# Setup Sui
WORKDIR /sui
RUN curl --location -o sui.tgz https://github.com/MystenLabs/sui/releases/download/testnet-v1.17.0/sui-testnet-v1.17.0-ubuntu-x86_64.tgz
RUN tar zxvf sui.tgz
RUN mv target/release/* .
RUN mv external-crates/move/target/release/* .

RUN mv sui-ubuntu-x86_64 sui
RUN mv sui-test-validator-ubuntu-x86_64 sui-test-validator
RUN mv move-analyzer-ubuntu-x86_64 move-analyzer

# RUN git clone --depth 1 --branch mainnet-v1.16.2 https://github.com/MystenLabs/sui.git
# RUN cargo install --locked --git https://github.com/MystenLabs/sui.git --tag mainnet-v1.16.2 sui

WORKDIR /
ENV PATH $PATH:/sui
RUN sui help
ENV RUST_LOG="error" 

ENTRYPOINT [ "bash" ]