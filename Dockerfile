FROM scratch

WORKDIR /app

COPY static /app/static
COPY server/coursedb.json /app/server/
COPY schedulecreator-backend /app

ENTRYPOINT ["/app/schedulecreator-backend"]
CMD [""]

EXPOSE 8080
