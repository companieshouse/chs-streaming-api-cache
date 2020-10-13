chs-streaming-api-cache
======================

Streaming api cache service. Consumes deltas from chs-streaming-api-backend and caches the entries in Redis.

The integration tests use the following environment variables

- REDIS_URL
- CACHE_EXPIRY_IN_SECONDS
