# goggler

Website screenshot microservice written in Golang that caches screenshots on disk or S3 storage

## Compiling and Running

Compile and run with:

```bash
go build github.com/mkalus/goggler/cmd/goggler && GOGGLER_DEBUG=1 ./goggler
```

Point your Browser to the URL: `http://localhost:8080/?url=https%3A%2F%2Fduckduckgo.com%2F&wait=2000`

**Note:** Goggler uses ChromeDP and needs Chrome or Chromium to work. This means it has to be able to
access and start Chrome or Chromium somehow (Chrome in path variable or the like).

## Parameters

Valid URI parameters:

* `url` URL to take screenshot from (required)
* `width` Width of image and viewport (might be scaled up or down, see below, default: 1920)
* `height` Height of image and viewport (might be scaled up or down, see below, default: 1024)
* `scale` Scale of image and viewport (final width and height are multiplied by this, default: 0.2)
* `quality` PNG image quality (default: 90)
* `wait` Time in ms to wait until screenshot is taken (to wait for rendering to finish, default: 2000 =  2 secs)
* `timeout` Timeout in ms after which request is cancelled (default: 60000 = 60 secs)
* `maxage` Maximum age of cache file in seconds (!) before it gets reloaded (set to 0 to never renew files, default: 2592000 = 30 days)
* `force` Force update (set to `1` or any other non-empty value)

Defaults may be changed by setting environmental variables when running goggle.

## Environmental variables

Changes defaults and sets some other elements:

* `GOGGLER_WIDTH` Set default image width
* `GOGGLER_HEIGHT` Set default image height
* `GOGGLER_SCALE` Set default image scale
* `GOGGLER_QUALITY` Set default image quality
* `GOGGLER_WAIT_FOR_IDLE` Wait for network idle in browser (default: `false`, set to any value to enable)
* `GOGGLER_WAIT` Set default wait time (ms) before rendering (default: `2000` aka 2 seconds)
* `GOGGLER_TIMEOUT` Set default timeout (ms) time (default: `60000` aka 60 seconds)
* `GOGGLER_MAXAGE` Set default max age (s) time (default: `2592000` aka 30 days)
* `GOGGLER_LISTEN` Set default listen address (default: `:8080`)
* `GOGGLER_DEBUG` Enable debugging log (default: `false`, set to any value to enable)
* `GOGGLER_CACHE` Type of cache (`local` or `s3`, default `local`)
* `GOGGLER_CACHE_CLEANUP_INTERVAL` Interval in seconds at which cleanup service is run to clean stale data (maxage has to be greater than 0, also: set to 0 to never clean up old files, s3 needs full days to work properly, default: 2592000 = 30 days)
* `GOGGLER_CACHE_LOCAL_PATH` Path to local cache (default: OS specific temp dir like `/tmp/goggler`)
* `GOGGLER_CACHE_S3_URL` Endpoint for S3 storage (e.g. `s3.amazonaws.com`)
* `GOGGLER_CACHE_S3_BUCKETNAME` Bucket name
* `GOGGLER_CACHE_S3_ACCESSKEY` Access Key ID for S3 storage
* `GOGGLER_CACHE_S3_SECRETKEY` Secret Access Key ID for S3 storage
* `GOGGLER_CACHE_S3_REGION` S3 region (might be skipped in Amazon, e.g. `us-east-1`)
* `GOGGLER_CACHE_S3_SKIPSSL` Set to any value to skip secure SSL/TLS connection
* `GOGGLER_CACHE_S3_CREATEBUCKET` Set to any value to create bucket if it does not exist

Examples:

```bash
# Basic test
GOGGLER_DEBUG=1 GOGGLER_LISTEN=127.0.0.1:9090 ./goggler
# Local storage
GOGGLER_CACHE_LOCAL_PATH=~/mydata ./goggler
# Waiting for network idle instead of fixed amount of time
GOGGLER_CACHE_LOCAL_PATH=~/mydata GOGGLER_WAIT_FOR_IDLE=1 ./goggler
# S3 storage
GOGGLER_DEBUG=1 GOGGLER_CACHE=s3 GOGGLER_CACHE_S3_BUCKETNAME=mytestbucket \
  GOGGLER_CACHE_S3_ACCESSKEY=_KEY_ GOGGLER_CACHE_S3_SECRETKEY=_KEY_ ./goggler
# Minio storage
GOGGLER_DEBUG=1 GOGGLER_CACHE=s3 GOGGLER_CACHE_S3_URL=127.0.0.1:9000 \
  GOGGLER_CACHE_S3_BUCKETNAME=test GOGGLER_CACHE_S3_ACCESSKEY=minioadmin \
  GOGGLER_CACHE_S3_SECRETKEY=minioadmin GOGGLER_CACHE_S3_SKIPSSL=1 \
  GOGGLER_CACHE_S3_CREATEBUCKET=1 ./goggler
```

## Docker

There is a Docker container of goggler including a headless version of Chromium. Try it using:

```bash
docker run --rm -p8080:8080 --init ronix/goggler
```

**Important:** The `--init` option is needed to get rid of zombie processes that will spawn if you run the container.
This is due Chrome creating new processes within the container. I have not found a way to tackle this, but since it is
more a feature than a bug, I assume we can live with this.

Full example with persistent volume:

```bash
docker run -d -p8080:8080 -v /tmp/goggler:/tmp/goggler \
  --name goggler --init ronix/goggler
```

Full example with S3 storage:

```bash
docker run -d -p8080:8080 --name goggler -e "GOGGLER_CACHE=s3" \
  -e "GOGGLER_CACHE_S3_BUCKETNAME=mytestbucket" \
  -e "GOGGLER_CACHE_S3_ACCESSKEY=_KEY_" \
  -e "GOGGLER_CACHE_S3_SECRETKEY=_KEY_" \
  --init ronix/goggler
```

Building Docker image from this source:

```bash
docker build -t goggler .
```
