FROM debian:9-slim

ADD /bin/imx /usr/local/bin/
ADD /scripts/docker-entrypoint.sh /

ENV REGISTRY 127.0.0.1:2379

WORKDIR /

EXPOSE 8080 9321

# see https://github.com/golang/go/commit/9dee7771f561cf6aee081c0af6658cc81fac3918
RUN echo 'hosts: files mdns4_minimal [NOTFOUND=return] dns mdns4' >> /etc/nsswitch.conf
RUN chmod +x docker-entrypoint.sh

CMD ["/docker-entrypoint.sh"]
