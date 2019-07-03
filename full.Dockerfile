FROM scratch

LABEL Author="Yoann Le Garff (le-garff-yoann) <pe.weeble@yahoo.fr>"

COPY dist/backend/plateau /
COPY vue/plateau/index.html /public/
COPY vue/plateau/dist /public/dist

ENTRYPOINT [ "/plateau" ]
