# Creamy Shortener

<a href="https://hub.docker.com/r/albinodrought/creamy-shortener">
  <img alt="albinodrought/creamy-shortener Docker Pulls" src="https://img.shields.io/docker/pulls/albinodrought/creamy-shortener">
</a>
<a href="https://github.com/AlbinoDrought/creamy-shortener/blob/master/LICENSE">
  <img alt="AGPL-3.0 License" src="https://img.shields.io/github/license/AlbinoDrought/creamy-shortener">
</a>

Barebones link "shortener" using multihashes: identical URLs will always shorten to the same shortened links.

## Building

### Without Docker

```
go get -d -v
go build
```

### With Docker

`docker build -t albinodrought/creamy-shortener .`

## Running

```
CREAMY_APP_URL="https://shortener.r.albinodrought.com" \
CREAMY_DATA_PATH=/data \
CREAMY_HTTP_PORT=80 \
CREAMY_POPULATED_HOSTS=localhost,google.com \
./creamy-shortener
```

- `CREAMY_APP_URL`: the externally-accessible URL this instance can be reached at

- `CREAMY_DATA_PATH`: the path to persist all data, defaults to `./data`

- `CREAMY_HTTP_PORT`: port to listen on, defaults to `3000`

- `CREAMY_POPULATED_HOSTS`: hosts (including port) to allow link shortening for, defaults to `localhost`

## Shortening

```
curl -X POST -d "link=https://example.com/foo?bar" "https://shortener.r.albinodrought.com/shorten"
```

Output:

`https://shortener.r.albinodrought.com/l/121b68747470733a2f2f6578616d706c652e636f6d2f666f6f3f626172`
