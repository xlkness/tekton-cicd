FROM ubuntu:20.04

RUN apt-get update && \
    apt-get install apt-utils -y && \
    apt-get install wget -y && \
    apt-get install libncurses5 -y && \
    apt-get install ansible -y && \
    apt-get install sshpass -y && \
    apt-get install git -y && \
    apt-get install subversion -y && \
    apt-get install sshpass -y && \
    wget http://security.ubuntu.com/ubuntu/pool/main/o/openssl/libssl1.0.0_1.0.2g-1ubuntu4.20_amd64.deb && \
    dpkg -i libssl1.0.0_1.0.2g-1ubuntu4.20_amd64.deb

RUN apt-get install curl -y

RUN mkdir /data && mkdir /cert

COPY ca.crt /cert
COPY harbor.joynova.com.crt /cert
COPY harbor.joynova.com.key /cert

COPY .kube /root/.kube
RUN chmod g-r /root/.kube/config
RUN chmod o-r /root/.kube/config

RUN mkdir -p /root/.local/share/helm/plugins/cm-push
COPY cm-push /root/.local/share/helm/plugins/cm-push
COPY helm /bin
COPY kubectl /bin
COPY mysql /bin
COPY kubectl-up /bin
COPY broadcast_reply /bin
COPY flyway-6.5.4/ /etc/flyway
COPY flyway /bin
COPY dbdiff /bin

WORKDIR /data

ENTRYPOINT [""]