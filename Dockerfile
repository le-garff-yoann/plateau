FROM scratch

LABEL Author="Yoann Le Garff (le-garff-yoann) <pe.weeble@yahoo.fr>"

COPY dist/backend/plateau /

ENTRYPOINT [ "/plateau" ]
