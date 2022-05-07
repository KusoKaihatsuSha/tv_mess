FROM golang:1.18.0-alpine as builder
WORKDIR /tvmess/
COPY *.go *.mod ./
RUN go get github.com/boltdb/bolt
RUN go get github.com/kkdai/youtube/v2
RUN go get github.com/google/uuid
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /bin/appimage

FROM alpine:latest
COPY --from=builder /bin/appimage /bin/appimage
RUN apk update
RUN apk add ffmpeg

ENTRYPOINT ["/bin/appimage"]
CMD [ "/bin/appimage" ]