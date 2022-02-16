FROM golang:alpine AS build
RUN apk --no-cache add gcc g++ make git
WORKDIR /go/src/app
COPY main .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./bin/web-app ./main.go

FROM alpine:3.13
RUN apk --no-cache add ca-certificates
WORKDIR /usr/bin
COPY --from=build /go/src/app/bin /go/bin
COPY --from=build /go/src/app/templateIndexHtml /var/www/build/templateIndexHtml
COPY --from=build /go/src/app/unsolved-std.png /var/www/img/unsolved-std.png
EXPOSE 80
ENTRYPOINT /go/bin/web-app --port 80
