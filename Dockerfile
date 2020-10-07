FROM golang:alpine as build

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN apk update && \
	apk upgrade && \
	apk add gcc bash && \
	rm -rf /var/cache/apk/*

WORKDIR /hubbot
COPY ./ /hubbot/
RUN go build -o bin/hubbot ./cmd

FROM amd64/alpine:3

EXPOSE 8080

LABEL maintainer="webhippie united mail@webhippie.de" \
  org.label-schema.name="webhippie hubbot" \
  org.label-schema.vendor="webhippie united" \
  org.label-schema.schema-version="0.1.0"

ENTRYPOINT ["/usr/bin/hubbot"]

COPY --from=build /hubbot/bin/hubbot /usr/bin/hubbot