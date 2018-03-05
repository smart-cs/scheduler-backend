FROM scratch

WORKDIR /app

COPY static /app/static
COPY server/coursedb.json /app/server/
COPY scheduler-backend /app

ENTRYPOINT ["/app/scheduler-backend"]
CMD [""]

EXPOSE 8080
