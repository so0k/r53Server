FROM golang:alpine as builder

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN	apk add --no-cache \
	ca-certificates

COPY . /go/src/github.com/so0k/r53Server

RUN set -x \
	&& apk add --no-cache --virtual .build-deps \
		git \
		gcc \
		libc-dev \
		libgcc \
		make \
	&& cd /go/src/github.com/so0k/r53Server \
	&& make static \
	&& mv r53Server /usr/bin/r53Server \
	&& apk del .build-deps \
	&& rm -rf /go \
	&& echo "Build complete."

FROM scratch

COPY --from=builder /usr/bin/r53Server /usr/bin/r53Server
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs

COPY static static
COPY templates templates

ENTRYPOINT [ "r53Server" ]
CMD [ "--help" ]
