FROM scratch

COPY server .

USER 1000

ENTRYPOINT ["/server"]
