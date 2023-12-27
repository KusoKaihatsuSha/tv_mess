FROM golang:1.21.0-alpine as builder
WORKDIR /tvmess/
COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /bin/appimage

FROM alpine:latest
COPY --from=builder /bin/appimage /bin/appimage
RUN apk update && apk add ffmpeg

ENTRYPOINT ["/bin/appimage"]
CMD [ "/bin/appimage" ]