FROM alpine

MAINTAINER Bo Yang <yangbo1010@gmail.com>

ADD build/bin/webservice-linux-amd64 /usr/bin/

ADD config.yaml /usr/bin/

ADD certs/ /usr/bin/certs/

EXPOSE 8087 8088

WORKDIR /usr/bin/

ENTRYPOINT ["/usr/bin/webservice"]
