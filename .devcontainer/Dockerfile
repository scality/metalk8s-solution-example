FROM mcr.microsoft.com/devcontainers/go:1.24-trixie

RUN export DEBIAN_FRONTEND=noninteractive && \
    apt-get update && \
    apt-get install --no-install-recommends -y \
        bash-completion \
        build-essential \
        curl \
        docker.io \
        genisoimage \
        gh \
        git \
        isomd5sum \
        libgpgme-dev \
        make \
        shellcheck \
        vim \
        yq \
        && \
    apt-get clean

USER vscode

ARG OPERATOR_SDK_VERSION=1.42.0

RUN curl -sSfLO https://github.com/operator-framework/operator-sdk/releases/download/v${OPERATOR_SDK_VERSION}/operator-sdk_linux_amd64 && \
    chmod +x operator-sdk_linux_amd64 && \
    sudo mv operator-sdk_linux_amd64 /usr/local/bin/operator-sdk

ARG GOLANGCI_VERSION=2.5.0

RUN curl -sSfL  https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sudo sh -s -- v${GOLANGCI_VERSION}

ARG KIND_VERSION=0.31.0

RUN curl -sSfLO https://kind.sigs.k8s.io/dl/v${KIND_VERSION}/kind-linux-amd64 && \
    chmod +x kind-linux-amd64 && \
    sudo mv kind-linux-amd64 /usr/local/bin/kind

ARG KUBECTL_VERSION=1.35.2

RUN curl -sSfLO https://dl.k8s.io/release/v${KUBECTL_VERSION}/bin/linux/amd64/kubectl && \
    chmod +x kubectl && \
    sudo mv kubectl /usr/local/bin/kubectl

ARG KUSTOMIZE_VERSION=5.8.1

RUN curl -sSfL -o /tmp/kustomize.tar.gz https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2Fv${KUSTOMIZE_VERSION}/kustomize_v${KUSTOMIZE_VERSION}_linux_amd64.tar.gz && \
    sudo tar -xzf /tmp/kustomize.tar.gz -C /usr/local/bin/ kustomize && \
    rm -rf /tmp/kustomize.tar.gz

ARG SKOPEO_VERSION=1.22.0

RUN curl -sSfL -o /tmp/skopeo.tar.gz https://github.com/containers/skopeo/archive/refs/tags/v${SKOPEO_VERSION}.tar.gz && \
    tar -zxvf /tmp/skopeo.tar.gz -C /tmp && \
    cd /tmp/skopeo-${SKOPEO_VERSION} && \
    sudo PATH=$PATH make bin/skopeo && \
    sudo mv bin/skopeo /usr/local/bin/ && \
    cd && \
    rm -rf /tmp/skopeo.tar.gz /tmp/skopeo-${SKOPEO_VERSION}
