FROM golang:alpine as build

COPY ./ /hubbot/
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN apk update && \
	apk upgrade && \
	apk add make gcc bash && \
	rm -rf /var/cache/apk/*

WORKDIR /hubbot/src
RUN make clean generate build


FROM amd64/alpine:3

LABEL maintainer="webhippie united" \
  org.label-schema.name="webhippie hubbot" \
  org.label-schema.vendor="webhippie united" \
  org.label-schema.schema-version="0.1.0"

ENTRYPOINT ["/usr/bin/hubbot"]

COPY --from=build /hubbot/bin/hubbot /usr/bin/hubbot