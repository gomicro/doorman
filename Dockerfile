FROM golang:1.12-alpine
MAINTAINER dev@gomicro.io

ADD . /go/src/github.com/gomicro/doorman
ADD https://github.com/gomicro/probe/releases/download/v0.0.3/probe_0.0.3_linux_amd64.tar.gz /probe.tar.gz
RUN tar xvf /probe.tar.gz -C /

WORKDIR /go/src/github.com/gomicro/doorman

RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build -o /service .

FROM scratch
MAINTAINER dev@gomicro.io

COPY --from=0 /service /service
COPY --from=0 /probe /probe

HEALTHCHECK --interval=5s --timeout=30s --retries=3 CMD ["/probe", "http://localhost:4567/v1/status"]

EXPOSE 4567

CMD ["/service"]
