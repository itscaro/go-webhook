docker build -t build . -f Dockerfile.redist
docker create --name build build
docker cp build:/root/webhook .
docker cp build:/root/webhook-darwin .
docker cp build:/root/webhook-armhf .
docker rm -f build
