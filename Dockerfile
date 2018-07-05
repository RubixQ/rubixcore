FROM golang:latest
LABEL author="Edward Pie"
ENV RUBIXCORE_PORT=5000
ENV RUBIXCORE_APP_ENV=development
ENV RUBIXCORE_MONGO_DSN=db
ENV SRC_DIR=/go/src/github.com/rubixq/rubixcore
ADD . ${SRC_DIR}
RUN go get -u github.com/golang/dep/cmd/dep
WORKDIR ${SRC_DIR}
RUN dep ensure -v
RUN go build .
ENTRYPOINT [ "./rubixcore" ]