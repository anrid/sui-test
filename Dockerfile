FROM ubuntu:22.04

RUN apt-get update && apt-get -y upgrade
RUN apt-get -y install curl git

RUN echo "alias ll='ls -ltr'" > /root/.bashrc
RUN echo "export PATH=$PATH:/sui" >> /root/.bashrc

RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y

WORKDIR /sui
RUN curl --location -o sui.tgz https://github.com/MystenLabs/sui/releases/download/mainnet-v1.16.2/sui-mainnet-v1.16.2-ubuntu-x86_64.tgz
RUN tar zxvf sui.tgz
RUN mv target/release/* .
RUN mv external-crates/move/target/release/* .
RUN mv sui-ubuntu-x86_64 sui
RUN mv sui-test-validator-ubuntu-x86_64 sui-test-validator

# RUN git clone --depth 1 --branch mainnet-v1.16.2 https://github.com/MystenLabs/sui.git
# RUN cargo install --locked --git https://github.com/MystenLabs/sui.git --tag mainnet-v1.16.2 sui

WORKDIR /
ENV PATH=$PATH:/sui
RUN sui help

RUN apt-get -y install libpq-dev jq
ENV RUST_LOG="error" 

ENTRYPOINT [ "sui-test-validator" ]