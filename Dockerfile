FROM alpine

RUN apk add bash
RUN apk add -U tzdata
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

WORKDIR /opt/tcpProxy

COPY ./TcpProxy_linux64 /opt/tcpProxy/TcpProxy_linux64
COPY ./config.json /opt/tcpProxy/config.json
COPY ./app.yml /opt/tcpProxy/app.yml
COPY ./html/static/js/* /opt/tcpProxy/html/static/js/
COPY ./html/static/css/* /opt/tcpProxy/html/static/css/
COPY ./html/static/img/* /opt/tcpProxy/html/static/img/
COPY ./html/index.html /opt/tcpProxy/html/index.html
COPY ./certs /opt/tcpProxy/certs/

EXPOSE 18081

CMD "./TcpProxy_linux64"