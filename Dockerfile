ARG GOLANG_VERSION
ARG ALPINE_VERSION
FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} as builder

RUN apk --no-cache --virtual .build-deps add make gcc musl-dev binutils-gold

# go.mod's krakend-jose/krakend-cel/krakend-opencensus/krakend-martian
# replace directives point at ../krakend-<name> (see REBASE_PLAYBOOK.md --
# their own module identity is inherited from upstream and doesn't match
# where they're actually hosted, so a version-pinned replace can't be used).
# The build context must therefore be the parent directory containing this
# repo and all four sibling repos as checkouts, not just this repo alone.
COPY krakend-ce /app
COPY krakend-jose /krakend-jose
COPY krakend-cel /krakend-cel
COPY krakend-opencensus /krakend-opencensus
COPY krakend-martian /krakend-martian
WORKDIR /app

RUN make build


FROM alpine:${ALPINE_VERSION}

LABEL maintainer="community@krakend.io"

RUN apk upgrade --no-cache --no-interactive && apk add --no-cache ca-certificates tzdata && \
    adduser -u 1000 -S -D -H krakend && \
    mkdir /etc/krakend && \
    echo '{ "version": 3 }' > /etc/krakend/krakend.json

COPY --from=builder /app/krakend /usr/bin/krakend

USER 1000

WORKDIR /etc/krakend

ENTRYPOINT [ "/usr/bin/krakend" ]
CMD [ "run", "-c", "/etc/krakend/krakend.json" ]

EXPOSE 8000 8090
