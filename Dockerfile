FROM golang:1.12.1

LABEL autodelete=true
ARG GO_TAGS

COPY . $GOPATH/src/plateau/
RUN \
    cd $GOPATH/src/plateau && \
    CGO_ENABLED=0 go build -tags="${GO_TAGS}" -o /tmp/plateau

FROM node:11.3

LABEL autodelete=true

COPY vue/plateau/ /tmp/plateau/
RUN \
    cd /tmp/plateau && \
    npm install && \
    npm run build

FROM scratch

LABEL Author="Yoann Le Garff (le-garff-yoann) <pe.weeble@yahoo.fr>"

COPY --from=0 /tmp/plateau /
COPY --from=1 /tmp/plateau/dist /public

ENTRYPOINT [ "/plateau" ]
