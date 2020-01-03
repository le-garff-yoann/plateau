FROM node:13.5

LABEL autodelete=true

COPY vue/plateau/ /tmp/plateau/
RUN \
    cd /tmp/plateau && \
    npm install && \
    npm run build

FROM nginx:1.17

LABEL Author="Yoann Le Garff (le-garff-yoann) <pe.weeble@yahoo.fr>"

COPY docker/nginx/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=0 /tmp/plateau/dist /usr/share/nginx/html
