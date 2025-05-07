FROM dpage/pgadmin4

ENV PGADMIN_DEFAULT_EMAIL="admin@pgadmin.com" \
    PGADMIN_DEFAULT_PASSWORD="admin" \
    PGADMIN_LISTEN_PORT="80"

COPY scripts/pg-admin-scripts /scripts
RUN chmod +x /scripts/auto_connect.sh

WORKDIR /scripts

ENTRYPOINT ["./auto_connect.sh"]