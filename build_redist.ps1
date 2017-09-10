(rm -Recurse -Force -ErrorAction SilentlyContinue ./dist) -or (mkdir ./dist)
docker build -t build . -f Dockerfile.redist
docker create --name build build
docker cp build:/root/build/webhook_linux_amd64.tar.gz ./dist/
docker cp build:/root/build/webhook_windows_amd64.tar.gz ./dist/
docker cp build:/root/build/webhook_darwin_amd64.tar.gz ./dist/
docker cp build:/root/build/webhook_armhf.tar.gz ./dist/
docker rm -f build
