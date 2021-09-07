FROM centos:7

ARG BUILDBOT_VERSION=2.0.1
ARG GO_VERSION=1.13.10
ARG OPERATOR_SDK_VERSION=v0.17.0

ENV LANG=en_US.utf8

WORKDIR /home/eve/workspace

RUN yum install -y --setopt=skip_missing_names_on_install=False \
    epel-release \
    && yum-config-manager --add-repo \
    https://download.docker.com/linux/centos/docker-ce.repo \
    && yum install -y --setopt=skip_missing_names_on_install=False \
    sudo \
    gcc \
    hardlink \
    make \
    genisoimage \
    python36 \
    python36-devel \
    python36-pip \
    skopeo \
    yum-utils \
    docker-ce-cli-18.09.6 \
    git \
    && adduser -u 1042 --home /home/eve eve --groups docker \
    && chown eve:eve /home/eve \
    && python3.6 -m pip install buildbot-worker==${BUILDBOT_VERSION} pyyaml

# Add eve to sudoers.
RUN echo "eve ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers.d/eve

# Install Go tooling to build the Kubernetes Operators.
ENV GOROOT /usr/local/go
ENV GOPATH /home/eve/.golang
ENV GOCACHE /home/eve/.cache
ENV PATH $GOPATH/bin:$GOROOT/bin:$PATH

RUN curl -ORL https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz      \
    && tar xzvf go${GO_VERSION}.linux-amd64.tar.gz                             \
    && rm go${GO_VERSION}.linux-amd64.tar.gz                                   \
    && mv go /usr/local                                                        \
    && curl -RLo /usr/bin/operator-sdk                                         \
        https://github.com/operator-framework/operator-sdk/releases/download/${OPERATOR_SDK_VERSION}/operator-sdk-${OPERATOR_SDK_VERSION}-x86_64-linux-gnu \
    && chmod +x /usr/bin/operator-sdk

USER eve
