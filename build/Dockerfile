FROM registry.access.redhat.com/ubi7/ubi-minimal:7.6

RUN microdnf install unzip -y

ENV OPERATOR=/usr/local/bin/externalsecret-operator \
    USER_UID=1001 \
    USER_NAME=externalsecret-operator \
    ONEPASSWORD_CLI_VERSION=v0.5.6-003

USER root

# install operator binary
COPY build/_output/bin/externalsecret-operator ${OPERATOR}

COPY build/bin /usr/local/bin
RUN /usr/local/bin/user_setup

# install 1password binary
RUN cd /tmp; curl https://cache.agilebits.com/dist/1P/op/pkg/${ONEPASSWORD_CLI_VERSION}/op_linux_amd64_${ONEPASSWORD_CLI_VERSION}.zip -o op_linux_amd64_${ONEPASSWORD_CLI_VERSION}.zip; unzip op_linux_amd64_${ONEPASSWORD_CLI_VERSION}.zip; mv ./op /usr/local/bin/
RUN gpg --keyserver hkp://keys.gnupg.net --recv-keys 3FEF9748469ADBE15DA7CA80AC2D62742012EA22
RUN cd /tmp; gpg --verify /tmp/op.sig /usr/local/bin/op || (echo "ERROR: Incorrect GPG signature for 1password op binary." && exit 1)

USER ${USER_UID}

ENTRYPOINT ["/usr/local/bin/entrypoint"]
