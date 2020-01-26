# vim:set ft=dockerfile:
FROM golang:1.12

COPY . /go/src/github.com/webitel/storage
WORKDIR /go/src/github.com/webitel/storage/

RUN GOOS=linux go get -d ./...
RUN GOOS=linux go install
RUN GIT_COMMIT=$(git rev-list -1 HEAD) && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags "-X github.com/webitel/storage/model.BuildNumber=$GIT_COMMIT" -a -o engine .

FROM scratch

LABEL maintainer="Vitaly Kovalyshyn"

ENV WEBITEL_MAJOR 20.02
ENV WEBITEL_REPO_BASE https://github.com/webitel

WORKDIR /
COPY --from=0 /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /go/src/github.com/webitel/storage/i18n /i18n
COPY --from=0 /go/src/github.com/webitel/storage/storage /

ENTRYPOINT ["./storage"]
