FROM scratch

WORKDIR /app

COPY static /app/static
COPY database/coursedb.json /app/database/
COPY scheduler-backend /app

ENTRYPOINT ["/app/scheduler-backend"]
CMD [""]

EXPOSE 8080
