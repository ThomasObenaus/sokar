FROM ubuntu:19.10
LABEL maintainer="Thomas Obenaus"

# This is the release of Nomad to pull in.
ARG NOMAD_VERSION=0.9.1
# This is the release of Consul to pull in.
ARG CONSUL_VERSION=1.5.0

# To build this image one can call
# docker build -t nomad-consul -f Dockerfile --build-arg NOMAD_VERSION=0.9.1 --build-arg CONSUL_VERSION=1.5.0 .

# Copy needed install scripts and configuration files
COPY install-nomad-with-consul /install-nomad-with-consul
COPY docker-entrypoint.sh docker-entrypoint.sh
COPY nomad.hcl /nomad.hcl
COPY consul.hcl /consul.hcl

# Install needed dependencies.
RUN apt-get update -y && apt-get install -y \
    dumb-init \
    iproute2 \
    ca-certificates \
    curl \
    lxc \
    iptables \
    unzip

# Install Docker from Docker Inc. repositories.
RUN curl -sSL https://get.docker.com/ | sh

# Install the magic wrapper which enableds docker in docker
# Taken from https://github.com/jpetazzo/dind
ADD ./wrapdocker /usr/local/bin/wrapdocker
RUN chmod +x /usr/local/bin/wrapdocker

# A volume is needed since each nomad job needs at least some disk space
VOLUME /var/lib/docker

# Install consul and nomad in the specified versions
RUN ./install-nomad-with-consul --version ${NOMAD_VERSION} --consul-version ${CONSUL_VERSION}

# Server RPC is used for internal RPC communication between agents and servers,
# and for inter-server traffic for the consensus algorithm (raft).
EXPOSE 4647

# Serf is used as the gossip protocol for cluster membership. Both TCP and UDP
# should be routable between the server nodes on this port.
EXPOSE 4648 4648/udp

# HTTP is the primary interface that applications use to interact with Nomad.
EXPOSE 4646

# Consul HTTP (ui) port (tcp)
EXPOSE 8500

ENTRYPOINT ["./docker-entrypoint.sh"]

CMD ["bin/sh"]

