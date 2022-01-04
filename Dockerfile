FROM golang:1.17.5-alpine3.15 as backend-builder
WORKDIR /src
COPY backend/go.mod backend/go.sum /src/
RUN go mod download
COPY backend /src
RUN GOOS=linux GOARCH=arm64 go build -tags dynamodb -o /usr/bin/backend /src/cmd/server

FROM node:17-alpine3.14 as frontend-builder
WORKDIR /src
COPY frontend /src
RUN yarn install && yarn build

FROM alpine:3.15 as reproxy-builder
ARG REPROXY_VERSION="v0.11.0"
ARG REPROXY_SHA256_SUM="35dd1cc3568533a0b6e1109e7ba630d60e2e39716eea28d3961c02f0feafee8e"
ADD "https://github.com/umputun/reproxy/releases/download/${REPROXY_VERSION}/reproxy_${REPROXY_VERSION}_linux_arm64.tar.gz" /tmp/reproxy.tar.gz
RUN SHELL="/bin/ash" \
    set -o pipefail \
    && sha256sum "/tmp/reproxy.tar.gz" \
    && echo "${REPROXY_SHA256_SUM}  /tmp/reproxy.tar.gz" | sha256sum -c \
    && tar -xzf /tmp/reproxy.tar.gz -C /usr/bin \
    && rm /tmp/reproxy.tar.gz

FROM alpine:3.15
ARG S6_OVERLAY_VERSION="v2.2.0.3"
ARG S6_OVERLAY_SHA256_SUM="a24ebad7b9844cf9a8de70a26795f577a2e682f78bee9da72cf4a1a7bfd5977e"
ADD "https://github.com/just-containers/s6-overlay/releases/download/${S6_OVERLAY_VERSION}/s6-overlay-aarch64-installer" /tmp/s6-overlay-installer
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
