FROM node:13.13 AS build

LABEL autodelete=true

COPY vue/plateau/ /tmp/plateau/

WORKDIR /tmp/plateau

RUN npm install && \
    npm run build

FROM nginx:1.29 AS base

LABEL org.opencontainers.image.authors="Yoann Le Garff (le-garff-yoann) <pe.weeble@yahoo.fr>"

COPY docker/nginx/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /tmp/plateau/dist /usr/share/nginx/html
