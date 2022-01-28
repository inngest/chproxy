---
title: Caching
category: Configuration
position: 203
---

`Chproxy` may be configured to cache responses. It is possible to create multiple
[cache-configs](https://github.com/Vertamedia/chproxy/blob/master/config/#cache_config) with various settings.
Response caching is enabled by assigning cache name to user. Multiple users may share the same cache.
Currently only `SELECT` responses are cached.
Caching is disabled for request with `no_cache=1` in query string.
Optional cache namespace may be passed in query string as `cache_namespace=aaaa`. This allows caching
distinct responses for the identical query under distinct cache namespaces. Additionally,
an instant cache flush may be built on top of cache namespaces - just switch to new namespace in order
to flush the cache.


