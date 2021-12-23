FROM golang:1.17.5-alpine3.15 as backend-builder
WORKDIR /src
COPY backend/go.mod backend/go.sum /src/
RUN go mod download
COPY backend /src
RUN GOOS=linux GOARCH=arm64 go build -o /usr/bin/backend /src/cmd/server

FROM node:17-alpine3.14 as frontend-builder
WORKDIR /src
COPY frontend /src
RUN yarn install && yarn build

FROM alpine:3.15 as reproxy-builder
ARG REPROXY_VERSION="v0.11.0"
ARG REPROXY_SHA256_SUM="e8de61f5eed761540f9eb371e98536a770b41befc7ac0798e778077d6bef5511"
ADD "https://github.com/umputun/reproxy/releases/download/${REPROXY_VERSION}/reproxy_${REPROXY_VERSION}_linux_arm.tar.gz" /tmp/reproxy.tar.gz
RUN SHELL="/bin/ash" \
    set -o pipefail \
    && sha256sum "/tmp/reproxy.tar.gz" \
    && echo "${REPROXY_SHA256_SUM}  /tmp/reproxy.tar.gz" | sha256sum -c \
    && tar -xzf /tmp/reproxy.tar.gz -C /usr/bin \
    && rm /tmp/reproxy.tar.gz

FROM alpine:3.15
ARG S6_OVERLAY_VERSION="v2.2.0.3"
ARG S6_OVERLAY_SHA256_SUM="7140eafc62720ecc43f81292c9bdd75bc7f4d0421518d42707e69f8e78e55088"
ADD "https://github.com/just-containers/s6-overlay/releases/download/${S6_OVERLAY_VERSION}/s6-overlay-amd64-installer" /tmp/s6-overlay-installer
RUN SHELL="/bin/ash" \
    set -o pipefail \
    && sha256sum "/tmp/s6-overlay-installer" \
    && echo "${S6_OVERLAY_SHA256_SUM}  /tmp/s6-overlay-installer" | sha256sum -c \
    && chmod +x /tmp/s6-overlay-installer
COPY --from=backend-builder /usr/bin/backend /usr/bin/backend
COPY --from=frontend-builder /src/dist /frontend/dist
COPY --from=reproxy-builder /usr/bin/reproxy /usr/bin/reproxy
RUN "/tmp/s6-overlay-installer" /
COPY s6 /etc
ENV S6_KILL_GRACETIME=0
ENTRYPOINT [ "/init" ]
EXPOSE 80
