FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /src/devops/ss-pods-reader/

COPY . .

RUN go get -d -v

RUN go build -o /go/bin/ss-pods-reader


FROM scratch

ARG HOST
ARG TOKEN

WORKDIR /app

COPY --from=builder /go/bin/ss-pods-reader /app
COPY --from=builder /src/devops/ss-pods-reader/.env /app
COPY --from=builder /src/devops/ss-pods-reader/pods-running-tpl.html /app
ENTRYPOINT ["./ss-pods-reader"]

EXPOSE 8080
