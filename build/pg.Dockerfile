FROM postgres:16.2

RUN mkdir -p /usr/share/postgresql/16/tsearch_data
RUN chmod 777 /usr/share/postgresql/16/tsearch_data