# Stage 1: Use image builder to build the target binaries
# Copyright Contributors to the Open Cluster Management project

FROM registry.ci.openshift.org/stolostron/builder:go1.23-linux AS builder

ENV COMPONENT=governance-policy-propagator
ENV REPO_PATH=/go/src/github.com/stolostron/${COMPONENT}
WORKDIR ${REPO_PATH}
COPY . .
RUN make build

# Stage 2: Copy the binaries from the image builder to the base image
FROM registry.access.redhat.com/ubi9/ubi-minimal:latest

ENV COMPONENT=governance-policy-propagator
ENV REPO_PATH=/go/src/github.com/stolostron/${COMPONENT}
ENV OPERATOR=/usr/local/bin/${COMPONENT} \
    USER_UID=1001 \
    USER_NAME=${COMPONENT}

# install operator binary
COPY --from=builder ${REPO_PATH}/build/_output/bin/${COMPONENT} ${OPERATOR}

COPY --from=builder ${REPO_PATH}/build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}

LABEL name="rhacm2/governance-policy-propagator-rhel9"
LABEL summary="Propagate Policies to managed cluster namespaces"
LABEL description="Propagate Policies from user namespaces to managed cluster namespaces as defined by the PlacementBinding and bound Placement."
LABEL io.k8s.display-name="Governance policy propagator"
LABEL io.k8s.description="Propagate Policies from user namespaces to managed cluster namespaces as defined by the PlacementBinding and bound Placement."
LABEL com.redhat.component="acm-governance-policy-propagator-container"
LABEL io.openshift.tags="data,images"
