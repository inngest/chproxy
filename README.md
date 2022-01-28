[![Go Report Card](https://goreportcard.com/badge/github.com/Vertamedia/chproxy)](https://goreportcard.com/report/github.com/Vertamedia/chproxy)
[![Build Status](https://travis-ci.org/Vertamedia/chproxy.svg?branch=master)](https://travis-ci.org/Vertamedia/chproxy?branch=master)
[![Coverage](https://img.shields.io/badge/gocover.io-75.7%25-green.svg)](http://gocover.io/github.com/Vertamedia/chproxy?version=1.9)

# chproxy

Chproxy is an http proxy and load balancer for [ClickHouse](https://ClickHouse.yandex) database. 

Full documentation is available on [chproxy.org](https://www.chproxy.org/)

It provides the following features:

- May proxy requests to multiple distinct `ClickHouse` clusters depending on the input user. For instance, requests from `appserver` user may go to `stats-raw` cluster, while requests from `reportserver` user may go to `stats-aggregate` cluster.
- May map input users to per-cluster users. This prevents from exposing real usernames and passwords used in `ClickHouse` clusters. Additionally this allows mapping multiple distinct input users to a single `ClickHouse` user.
- May accept incoming requests via HTTP and HTTPS.
- May limit HTTP and HTTPS access by IP/IP-mask lists.
- May limit per-user access by IP/IP-mask lists.
- May limit per-user query duration. Timed out or canceled queries are forcibly killed
  via [KILL QUERY](http://clickhouse-docs.readthedocs.io/en/latest/query_language/queries.html#kill-query).
- May limit per-user requests rate.
- May limit per-user number of concurrent requests.
- All the limits may be independently set for each input user and for each per-cluster user.
- May delay request execution until it fits per-user limits.
- Per-user [response caching](#caching) may be configured.
- Response caches have built-in protection against [thundering herd](https://en.wikipedia.org/wiki/Cache_stampede) problem aka `dogpile effect`.
- Evenly spreads requests among replicas and nodes using `least loaded` + `round robin` technique.
- Monitors node health and prevents from sending requests to unhealthy nodes.
- Supports automatic HTTPS certificate issuing and renewal via [Let’s Encrypt](https://letsencrypt.org/).
- May proxy requests to each configured cluster via either HTTP or [HTTPS](https://github.com/yandex/ClickHouse/blob/96d1ab89da451911eb54eccf1017eb5f94068a34/dbms/src/Server/config.xml#L15).
- Prepends User-Agent request header with remote/local address and in/out usernames before proxying it to `ClickHouse`, so this info may be queried from [system.query_log.http_user_agent](https://github.com/yandex/ClickHouse/issues/847).
- Exposes various useful [metrics](#metrics) in [prometheus text format](https://prometheus.io/docs/instrumenting/exposition_formats/).
- Configuration may be updated without restart - just send `SIGHUP` signal to `chproxy` process.
- Easy to manage and run - just pass config file path to a single `chproxy` binary.
- Easy to [configure](https://github.com/Vertamedia/chproxy/blob/master/config/examples/simple.yml):
```yml
server:
  http:
    listen_addr: ":9090"
    allowed_networks: ["127.0.0.0/24"]

users:
  - name: "default"
    to_cluster: "default"
    to_user: "default"

# by default each cluster has `default` user which can be overridden by section `users`
clusters:
  - name: "default"
    nodes: ["127.0.0.1:8123"]

```

## How to install

### Precompiled binaries

Precompiled `chproxy` binaries are available [here](https://github.com/Vertamedia/chproxy/releases).
Just download the latest stable binary, unpack and run it with the desired [config](#configuration):

```
./chproxy -config=/path/to/config.yml
```

### Building from source

Chproxy is written in [Go](https://golang.org/). The easiest way to install it from sources is:

```
go get -u github.com/Vertamedia/chproxy
```

If you don't have Go installed on your system - follow [this guide](https://golang.org/doc/install).


