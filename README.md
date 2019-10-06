# Creamy Shortener

<a href="https://hub.docker.com/r/albinodrought/creamy-shortener">
  <img alt="albinodrought/creamy-shortener Docker Pulls" src="https://img.shields.io/docker/pulls/albinodrought/creamy-shortener">
</a>
<a href="https://github.com/AlbinoDrought/creamy-shortener/blob/master/LICENSE">
  <img alt="AGPL-3.0 License" src="https://img.shields.io/github/license/AlbinoDrought/creamy-shortener">
</a>

Barebones link "shortener" using [multihashes](https://github.com/multiformats/multicodec): identical URLs will always shorten to the same shortened links.


## Shortening

```
curl -X POST -d "link=https://example.com/foo?bar" "https://your.shortener/shorten"
```

Output:

`https://your.shortener/l/Qma3YMYZUNAY7Dp7UhtZfqKAfsLkHyF9jf1yFXjZbYjWqt`


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
CREAMY_APP_URL="https://your.shortener/" \
CREAMY_DATA_PATH=/data \
CREAMY_HASH_MODE=sha2-256 \
CREAMY_HTTP_PORT=80 \
CREAMY_POPULATED_HOSTS=localhost,example.com \
./creamy-shortener
```

- `CREAMY_APP_URL`: the externally-accessible URL this instance can be reached at, defaults to `http://localhost:3000/`

- `CREAMY_DATA_PATH`: the path to persist all data, defaults to `./data`

- `CREAMY_HASH_MODE`: [multihash mode](https://github.com/multiformats/multicodec/blob/8dd1bfb9a953da79c3ec5962a2a3dcb94e0cc376/table.csv) to use, defaults to `sha2-256`

- `CREAMY_HTTP_PORT`: port to listen on, defaults to `3000`

- `CREAMY_POPULATED_HOSTS`: hosts (including port) to allow link shortening for, defaults to `localhost`
