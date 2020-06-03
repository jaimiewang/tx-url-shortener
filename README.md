# tx-url-shortener
Simple high performance URL shortener microservice written in Go.

## Docker installation
```shell script
docker build -t tx-url-shortener:1.0 .
docker run -d -v $(pwd)/config.yml:/opt/tx-url-shortener/config.yml
              -v $(pwd)/db.sqlite3:/opt/tx-url-shortener/db.sqlite3
              --publish 8080:8080
              --name tx-url-shortener
              tx-url-shortener:1.0
```

## Endpoints
At present there are only two basic endpoints for viewing data about specified URL and
for shortening URLs.