# tx-url-shortener

Simple high performance URL shortener microservice written in Go.

## Running as Docker container
```shell script
docker build -t tx-url-shortener:1.0 .
docker run -d -v $(pwd)/config.yml:/opt/tx-url-shortener/config.yml
              -v $(pwd)/db.sqlite3:/opt/tx-url-shortener/db.sqlite3
              --publish 8080:8080
              --name tx-url-shortener
              tx-url-shortener:1.0
```

## Endpoints
Actually, there aren't too many endpoints, exists only two basic for viewing data about
and shortening URL.