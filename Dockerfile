FROM golang:latest
LABEL author="Edward Pie"
ENV RUBIXCORE_PORT=5000
ENV RUBIXCORE_APP_ENV=development
ENV RUBIXCORE_POSTGRES_DSN="host=db port=5432 user=postgres password=rub1xc0r3 dbname=postgres sslmode=disable"
ENV RUBIXCORE_REDIS_URL=redis
ENV RUBIXCORE_TICKET_RESET_INTERVAL=12
ENV RUBIXCORE_JWT_ISSUER=rubixcore
ENV RUBIXCORE_JWT_SECRET=rub1xc0r3s3cr3tp@554jwt5
ENV DEFAULT_ADMIN_USERNAME=admin
ENV DEFAULT_ADMIN_PASSWORD=p@554@dm1n
ENV SRC_DIR=/go/src/github.com/rubixq/rubixcore
ADD . ${SRC_DIR}
RUN go get -u github.com/golang/dep/cmd/dep
WORKDIR ${SRC_DIR}
RUN dep ensure -v
RUN go build -race .
ENTRYPOINT [ "./rubixcore" ]