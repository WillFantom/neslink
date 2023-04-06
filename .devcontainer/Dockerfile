FROM ghcr.io/willfantom/devcontainer-go:latest

RUN apt-get update -yq && \
    apt-get install -y --no-install-recommends \
    iperf3 \
    iproute2 \
    net-tools \
    openvswitch-common \
    openvswitch-switch \
    openvswitch-testcontroller && \
    rm -rf /var/lib/apt/lists/*
