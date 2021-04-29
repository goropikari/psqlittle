FROM golang:1.16

WORKDIR /app
COPY . /app/

RUN go mod tidy
RUN make && cp ./bin/repl /usr/bin/repl && rm -rf /app

CMD /usr/bin/repl
