docker build -t build . -f Dockerfile.redist
docker create --name build build
docker cp build:/root/build/webhook .
docker cp build:/root/build/webhook.exe .
docker cp build:/root/build/webhook-darwin .
docker cp build:/root/build/webhook-armhf .
docker rm -f build
