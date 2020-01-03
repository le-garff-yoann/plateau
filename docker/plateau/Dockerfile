FROM golang:1.12.1

LABEL autodelete=true
ARG GO_TAGS

COPY . $GOPATH/src/plateau/
RUN \
    cd $GOPATH/src/plateau && \
    CGO_ENABLED=0 go build -tags="${GO_TAGS}" -o /tmp/plateau

FROM scratch

LABEL Author="Yoann Le Garff (le-garff-yoann) <pe.weeble@yahoo.fr>"

COPY --from=0 /tmp/plateau /

ENTRYPOINT [ "/plateau" ]
