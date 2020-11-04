# chs-streaming-api-cache

## Contents

The Companies House Streaming Platform Cache consumes offsets from the Streaming Platform Backend service and caches the entries in Redis, and pushes these to connected users as an event stream.

## Requirements

The following services and applications are required to build and/or run chs-streaming-api-cache:

* AWS ECR
* Docker
* Redis
* Companies House Streaming Platform Backend (chs-streaming-api-backend)

You will need an HTTP client that supports server-sent events (e.g. cURL) to connect to the service and receive published offsets.

## Building and Running Locally

1. Login to AWS ECR.
2. Build the project by running `docker build` from the project directory.
3. Run the Docker image that has been built by running `docker run IMAGE_ID` from the command line, ensuring values have been specified for the environment variables (see Configuration) and that port 6001 is exposed.
4. Send a GET request using your HTTP client to /filings. A connection should be established and any offsets published to the stream-filing-history topic should appear in the response body. The offsets should also be cached in Redis for other consumers.

## Configuration

Variable|Description|Example|Mandatory|
--------|-----------|-------|---------|
REDIS_URL|The URL of the Redis cache|redis:6379|yes
STREAMING_BACKEND_URL|The URL of the CH Streaming Backend service|http://chs-streaming-api-backend:6000|yes
REDIS_POOL_SIZE|The number of connections in a Redis connection pool|10|yes
CACHE_EXPIRY_IN_SECONDS|The number of seconds before a offset cache entry expires|3600|yes
STREAM_BACKEND_FILINGS_PATH|The backend endpoint to stream filing history offsets|/filings|yes
STREAM_BACKEND_COMPANIES_PATH|The backend endpoint to stream filing history offsets|/companies?timepoint=2|yes
STREAM_BACKEND_INSOLVENCY_PATH|The backend endpoint to stream company insolvency offsets|/insolvency-cases|yes
STREAM_BACKEND_CHARGES_PATH|The backend endpoint to stream company charges offsets|/charges|yes
STREAM_BACKEND_OFFICERS_PATH|The backend endpoint to stream officer appointments offsets|/officers|yes
STREAM_BACKEND_PSCS_PATH|The backend endpoint to stream PSC offsets|/persons-with-significant-control|yes