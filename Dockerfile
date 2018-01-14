FROM scratch

COPY ["static/", "/static/"]
COPY ["schedulecreator-backend", "coursedb.json", "/"]
ENTRYPOINT ["/schedulecreator-backend"]
CMD [""]

EXPOSE 8080
