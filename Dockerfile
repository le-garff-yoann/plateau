ARG PACKAGING=simple

FROM scratch AS build_simple

FROM scratch AS build_full

ONBUILD COPY vue/plateau/dist /

FROM build_${PACKAGING}

LABEL Author="Yoann Le Garff (le-garff-yoann) <pe.weeble@yahoo.fr>"

COPY dist/backend/plateau /

ENTRYPOINT [ "/plateau" ]
