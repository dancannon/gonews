FROM mischief/docker-golang

ADD . /root/go/src/github.com/dancannon/gonews
# Fetch deps
WORKDIR /root/go/src/github.com/dancannon/gonews

RUN go get github.com/tools/godep
RUN godep go install

# Allow to mount the current version in the container
VOLUME /root/go/src/github.com/dancannon/gonews

EXPOSE 3000
ENV PORT 3000

CMD /root/go/bin/gonews
