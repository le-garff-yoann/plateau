FROM scratch

LABEL Author="Yoann Le Garff (le-garff-yoann) <pe.weeble@yahoo.fr>"

COPY dist/backend/plateau /
COPY vue/plateau/index.html /public/
COPY vue/plateau/dist /public/dist

ENV PLATEAU_LISTENER 8080

EXPOSE 8080

ENTRYPOINT [ "/plateau" ]

CMD [ \
    "--listen", "${PLATEAU_LISTENER}", \
    "--listen-static-dir", "/public/", \
    "--session-key", "${PLATEAU_SESSION_KEY}", \
    "--rethinkdb-address", "${PLATEAU_RETHINKDB_ADDRESS}", \
    "--rethinkdb-database", "${PLATEAU_RETHINKDB_DATABASE}", \
    "--rethinkdb-create-tables", \
]
