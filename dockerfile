FROM postgres:16

COPY ./init/create_schema.sql /docker-entrypoint-initdb.d/

EXPOSE 5432

