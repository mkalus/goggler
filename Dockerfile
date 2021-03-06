# Build image
FROM golang:1.15-alpine AS build
WORKDIR /go/src/
COPY . github.com/mkalus/goggler
RUN cd github.com/mkalus/goggler && CGO_ENABLED=0 go build github.com/mkalus/goggler/cmd/goggler

# Actual image containing headless chromium browser
FROM demisto/chromium:1.0.0.16237
COPY --from=build /go/src/github.com/mkalus/goggler/goggler /opt/google/chrome/goggler
EXPOSE 8080
VOLUME /tmp/goggler
ENTRYPOINT ["/opt/google/chrome/goggler"]

# Alternatively, we could run the container like this:
#FROM demisto/chromium:1.0.0.16237
#RUN rm /etc/apt/sources.list.d/google-chrome.list && apt update --allow-unauthenticated && apt install -y dumb-init
#COPY --from=build /go/src/github.com/mkalus/goggler/goggler /opt/google/chrome/goggler
#EXPOSE 8080
#VOLUME /tmp/goggler
#ENTRYPOINT ["dumb-init", "--"]
#CMD ["/opt/google/chrome/goggler"]