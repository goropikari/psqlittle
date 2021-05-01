FROM golang:1.16.3-buster AS builder

WORKDIR /app
COPY . /app/

RUN go mod tidy
RUN make

FROM debian:buster-slim
ENV DBMS_HOST 0.0.0.0
WORKDIR /root/
COPY --from=builder /app/bin/repl /bin/
COPY --from=builder /app/bin/server /bin/

EXPOSE 5432

CMD /bin/server
