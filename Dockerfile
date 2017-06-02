FROM node:slim
MAINTAINER Vitaly Kovalyshyn "v.kovalyshyn@webitel.com"

ENV VERSION
ENV WEBITEL_MAJOR 3.4
ENV WEBITEL_REPO_BASE https://github.com/webitel

COPY src /cdr

WORKDIR /cdr
RUN npm install

EXPOSE 10021 10023
ENTRYPOINT ["node", "server.js"]
