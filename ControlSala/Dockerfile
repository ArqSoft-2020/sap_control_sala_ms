# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:latest

WORKDIR /go/src/app


COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .
EXPOSE 3030
CMD ["./main"]


FROM mysql:latest
RUN sed -i "s/^user.*/user = root/g" /etc/mysql/my.cnf

RUN chown -R mysql /var/lib/mysql
RUN chgrp -R mysql /var/lib/mysql
