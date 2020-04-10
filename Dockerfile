FROM scratch
MAINTAINER dev@gomicro.io

ADD doorman doorman
ADD ext/probe probe

EXPOSE 4567

CMD ["/doorman"]
