FROM golang:1.15-alpine AS build-env
RUN mkdir /go/src/app
ADD . /go/src/app
WORKDIR /go/src/app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o n9e-server .

FROM alpine:3.11
WORKDIR /app
COPY --from=build-env /go/src/app/n9e-server .
ENTRYPOINT ["./n9e-server"]
EXPOSE 8000