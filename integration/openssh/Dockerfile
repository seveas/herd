FROM alpine:latest

RUN apk --no-cache add openssh-server

EXPOSE 22

RUN ssh-keygen -t rsa -f /etc/ssh/ssh_host_rsa_key -N '' && \
    ssh-keygen -t ecdsa -f /etc/ssh/ssh_host_ecdsa_key -N '' && \
    ssh-keygen -t ed25519 -f /etc/ssh/ssh_host_ed25519_key -N '' && \
    mkdir /root/.ssh

ADD user.key.pub /root/.ssh/authorized_keys

ENTRYPOINT ["/usr/sbin/sshd", "-D", "-e", "-o", "PermitRootLogin=yes"]
