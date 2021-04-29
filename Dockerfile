FROM golang:1.16

ENV DBMS_HOST 0.0.0.0
WORKDIR /app
COPY . /app/

RUN go mod tidy
RUN make && cp ./bin/* /usr/bin/ && rm -rf /app

CMD /usr/bin/server
