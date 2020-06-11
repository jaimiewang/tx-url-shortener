# tx-url-shortener
Simple high performance URL shortener microservice written in Go.

## Installation
```shell script
docker build -t tx-url-shortener:1.0.4 .
docker run -d -v $(pwd)/config.yml:/opt/tx-url-shortener/config.yml \
              --publish 8080:8080 \
              --name tx-url-shortener \
              tx-url-shortener:1.0.4
```

## Generate new API key
```shell script
docker exec tx-url-shortener /opt/tx-url-shortener/bin/tx-url-shortener -generate-api-key
```

## Endpoints
### Shorten new URL
**Request**:
```shell script
curl -H "Authorization: Bearer <your-api-key>" \
     -H "Content-Type: application/json" \
     -X PUT \
     -d '{"url": "https://google.pl/"}' \
      http://localhost:8080/api/urls
```
**Response**:
```json
{
  "ip_address": "172.17.0.1",
  "views": 0,
  "code": "SPHTk",
  "created_at": 1591187797,
  "original": "https://google.pl/",
  "url": "http://localhost:8080/SPHTk"
}
```

### Get data about shortened URL
**Request**:
```shell script
curl -H "Authorization: Bearer <your-api-key>" \
     -X GET \
      http://localhost:8080/api/urls/<code>
```
**Response**:
```json
{
  "ip_address": "172.17.0.1",
  "views": 10,
  "code": "SPHTk",
  "created_at": 1591187797,
  "original": "https://google.pl/",
  "url": "http://localhost:8080/SPHTk"
}
```