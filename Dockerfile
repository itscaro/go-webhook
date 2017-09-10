# Building stage

FROM golang as builder

WORKDIR /go/src/github.com/itscaro/webhook
COPY . .

RUN go get -d -v

RUN GIT_COMMIT=$(git rev-list -1 HEAD) \
 && CGO_ENABLED=0 GOOS=linux go build --ldflags "-X main.GitCommit=${GIT_COMMIT} -X main.BuildDatetime=`date -u +%Y%m%d.%H%M`" -a -installsuffix cgo -o webhook . \
 && echo "Done"

# Packaging stage

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

ENV PORT=8080
ENV GIN_MODE=release
ENV UPNP_ENABLED=false
ENV UPNP_LOCAL_IP_RANGE=""

VOLUME ["./hook"]

EXPOSE 8080

COPY --from=builder /go/src/github.com/itscaro/webhook/webhook    .

CMD ["./webhook"]
