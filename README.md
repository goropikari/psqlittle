```
docker build -t dbms .
docker run -it -p 5432:5432 dbms  # server mode
psql -h 127.0.0.1 -p 5432  # connect dbms by using psql

docker run -it dbms repl
```
