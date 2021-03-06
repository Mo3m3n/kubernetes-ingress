# Copyright 2019 HAProxy Technologies LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang:1.16-alpine AS builder

RUN apk --no-cache add git openssh

COPY /go.mod /src/go.mod
COPY /go.sum /src/go.sum
RUN cd /src && go mod download

COPY / /src

RUN mkdir -p /var/run/vars && \
    cd /src && \
    git config --get remote.origin.url > /var/run/vars/GIT_REPO && \
    git rev-parse --short HEAD > /var/run/vars/GIT_HEAD_COMMIT && \
    git describe --abbrev=0 --tags > /var/run/vars/GIT_LAST_TAG && \
    git rev-parse --short $(cat /var/run/vars/GIT_LAST_TAG) > /var/run/vars/GIT_TAG_COMMIT && \
    git diff $(cat /var/run/vars/GIT_HEAD_COMMIT) $(cat /var/run/vars/GIT_TAG_COMMIT) --quiet > /var/run/vars/GIT_MODIFIED1 || echo '.dev' > /var/run/vars/GIT_MODIFIED1 && \
    git diff --quiet > /var/run/vars/GIT_MODIFIED2 || echo '.dirty' > /var/run/vars/GIT_MODIFIED2 && \
    cat /var/run/vars/GIT_MODIFIED1 /var/run/vars/GIT_MODIFIED2 | tr -d '\n' > /var/run/vars/GIT_MODIFIED && \
    date '+%Y-%m-%dT%H:%M:%S' > /var/run/vars/BUILD_DATE && \
    CGO_ENABLED=0 go build \
        -ldflags "-X main.GitRepo=$(cat /var/run/vars/GIT_REPO) -X main.GitTag=$(cat /var/run/vars/GIT_LAST_TAG) -X main.GitCommit=$(cat /var/run/vars/GIT_HEAD_COMMIT) -X main.GitDirty=$(cat /var/run/vars/GIT_MODIFIED) -X main.BuildTime=$(cat /var/run/vars/BUILD_DATE)" \
        -o fs/haproxy-ingress-controller .

FROM haproxytech/haproxy-alpine:2.4

COPY /fs/start.sh /
COPY /fs/usr/local/etc/haproxy/haproxy.cfg /usr/local/etc/haproxy/
COPY --from=builder /src/fs/haproxy-ingress-controller .

RUN apk --no-cache add socat openssl util-linux htop tzdata dumb-init && \
    rm -f /usr/local/bin/dataplaneapi /usr/bin/dataplaneapi && \
    ln -sf /usr/bin/dumb-init /dumb-init &&\
    chgrp -R haproxy /usr/local/etc/haproxy /run /var && \
    chmod -R g+w /usr/local/etc/haproxy /run /var

ENTRYPOINT ["/dumb-init", "--", "/start.sh"]
