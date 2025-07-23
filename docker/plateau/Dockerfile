FROM golang:1.23.1 AS build

LABEL autodelete=true
ARG GO_TAGS

COPY . $GOPATH/src/plateau/

WORKDIR $GOPATH/src/plateau

RUN CGO_ENABLED=0 go build -tags="${GO_TAGS}" -o /tmp/plateau

FROM scratch AS base

LABEL org.opencontainers.image.authors="Yoann Le Garff (le-garff-yoann) <pe.weeble@yahoo.fr>"

COPY --from=build /tmp/plateau /

ENTRYPOINT [ "/plateau" ]
