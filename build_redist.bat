(rm -Recurse -Force ./dist) -and (mkdir ./dist)
docker build -t build . -f Dockerfile.redist
docker create --name build build
docker cp build:/root/build/webhook ./dist/
docker cp build:/root/build/webhook.exe ./dist/
docker cp build:/root/build/webhook-darwin ./dist/
docker cp build:/root/build/webhook-armhf ./dist/
docker rm -f build
