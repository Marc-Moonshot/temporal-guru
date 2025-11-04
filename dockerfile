FROM postgres:16

# ENV POSTGRES_USER=admin \
#     POSTGRES_PASSWORD=temporal-guru-1337 \
#     POSTGRES_DB=mydb

# Copy initialization scripts if you have any SQL or shell setup scripts
# COPY ./initdb /docker-entrypoint-initdb.d/

# Expose the PostgreSQL default port
EXPOSE 5432

