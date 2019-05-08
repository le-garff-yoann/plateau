FROM scratch

LABEL Author="Yoann Le Garff (le-garff-yoann) <pe.weeble@yahoo.fr>"

COPY dist/backend/plateau /

ENV PLATEAU_LISTENER 8080

EXPOSE 8080

ENTRYPOINT [ "/plateau" ]

CMD [ \
    "--listen", "${PLATEAU_LISTENER}", \
    "--listen-session-key", "${PLATEAU_LISTENER_SESSION_KEY}", \
    "--pg-conn-str", "${PLATEAU_PG_CONNECTION_STRING}", \
]
