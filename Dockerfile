ARG PACKAGING=simple

FROM scratch AS build_simple

ONBUILD CMD [ "run", \
    "--listen", ":$LISTEN", \
    "--session-key", "$SESSION_KEY" \
]

FROM scratch AS build_full

ONBUILD CMD [ "run", \
    "--listen", ":$LISTEN", \
    "--session-key", "$SESSION_KEY", \
    "--listen-static-dir", "/public" \
]

ONBUILD COPY vue/plateau/dist /public

FROM build_${PACKAGING}

LABEL Author="Yoann Le Garff (le-garff-yoann) <pe.weeble@yahoo.fr>"

COPY dist/backend/plateau /

ENTRYPOINT [ "/plateau" ]
