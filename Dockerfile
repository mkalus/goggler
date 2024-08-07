# Build image
FROM docker.io/golang:1.22-alpine AS build
WORKDIR /go/src/
COPY . github.com/mkalus/goggler
RUN cd github.com/mkalus/goggler && CGO_ENABLED=0 go build github.com/mkalus/goggler/cmd/goggler

# Actual image containing headless chromium browser
FROM docker.io/demisto/chromium:126.0.6478.103218
COPY --from=build /go/src/github.com/mkalus/goggler/goggler /opt/google/chrome/goggler
EXPOSE 8080
VOLUME /tmp/goggler
ENTRYPOINT ["/opt/google/chrome/goggler"]

# Alternatively, we could run the container like this:
#FROM docker.io/demisto/chromium:1.0.0.56296
#RUN rm /etc/apt/sources.list.d/google-chrome.list && apt update --allow-unauthenticated && apt install -y dumb-init
#COPY --from=build /go/src/github.com/mkalus/goggler/goggler /opt/google/chrome/goggler
#EXPOSE 8080
#VOLUME /tmp/goggler
#ENTRYPOINT ["dumb-init", "--"]
#CMD ["/opt/google/chrome/goggler"]