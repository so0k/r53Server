FROM golang:alpine as builder

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN	apk add --no-cache \
	ca-certificates

COPY . /go/src/github.com/so0k/r53server

RUN set -x \
	&& apk add --no-cache --virtual .build-deps \
		git \
		gcc \
		libc-dev \
		libgcc \
		make \
	&& cd /go/src/github.com/so0k/r53server \
	&& make static \
	&& mv r53server /usr/bin/r53server \
	&& apk del .build-deps \
	&& rm -rf /go \
	&& mkdir /empty \
	&& echo "Build complete."

FROM scratch

COPY --from=builder /usr/bin/r53server /usr/bin/r53server
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs
COPY --from=builder /empty /tmp

COPY static static
COPY templates templates

ENTRYPOINT [ "r53server" ]
CMD [ "--help" ]
