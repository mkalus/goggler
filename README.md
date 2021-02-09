# goggler

Website screenshot microservice written in Golang that caches screenshots on disk or S3 storage

## Compiling and Running

Compile with:

```bash
go build github.com/mkalus/goggler/cmd/goggler
```

Point your Browser to the URL: `http://localhost:8080/?url=https%3A%2F%2Fduckduckgo.com%2F&wait=2000`

**Note:** Googler uses ChromeDP and needs Chrome or Chromium to work. This means it has to be able to
access and start Chrome or Chromium somehow (Chrome in path variable or the like).

## Docker

There is a Docker container of goggler including a headless version of Chromium. Try it using:

```bash
docker run --rm -p8080:8080 ronix/goggler
```

Building Docker image from this source:

```bash
docker build -t googler .
```
