# goggler

Website screenshot microservice written in Golang that caches screenshots on disk or S3 storage

## Compiling and Running

Compile and run with:

```bash
go build github.com/mkalus/goggler/cmd/goggler && GOGGLER_DEBUG=1 ./goggler
```

Point your Browser to the URL: `http://localhost:8080/?url=https%3A%2F%2Fduckduckgo.com%2F&wait=2000`

**Note:** Googler uses ChromeDP and needs Chrome or Chromium to work. This means it has to be able to
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
* `maxeage` Maximum age of cache file in seconds (!) before it gets reloaded (set to 0 to never renew files, default: 2592000 = 30 days)

Defaults may be changed by setting environmental variables when running goggle.

## Environmental variables

Changes defaults and sets some other elements:

* `GOGGLER_WIDTH` Set default image width
* `GOGGLER_HEIGHT` Set default image height
* `GOGGLER_SCALE` Set default image scale
* `GOGGLER_QUALITY` Set default image quality
* `GOGGLER_WAIT` Set default wait time
* `GOGGLER_TIMEOUT` Set default timeout time
* `GOGGLER_MAXAGE` Set default max age time
* `GOGGLER_LISTEN` Set default listen address (default: `:8080`)
* `GOGGLER_DEBUG` Enable debugging log
* `GOGGLER_CACHE` Type of cache (`local` or `s3`, default `local`)
* `GOGGLER_CACHE_LOCAL_PATH` Path to local cache (default: OS specific temp dir like `/tmp`)

Example:

```bash
GOGGLER_DEBUG=1 GOGGLER_LISTEN=127.0.0.1:9090 ./goggler
```

## Docker

There is a Docker container of goggler including a headless version of Chromium. Try it using:

```bash
docker run --rm -p8080:8080 ronix/goggler
```

Building Docker image from this source:

```bash
docker build -t googler .
```
