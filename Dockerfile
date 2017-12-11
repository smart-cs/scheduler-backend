FROM golang:latest

ADD . /go/src/github.com/nickwu241/schedulecreator-backend
RUN go install github.com/nickwu241/schedulecreator-backend

ENTRYPOINT ["/go/bin/schedulecreator-backend"]

EXPOSE 8080
