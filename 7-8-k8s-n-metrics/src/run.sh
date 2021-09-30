docker run \
  -v $(pwd):/var/loadtest \
  --net host \
  -it \
  --entrypoint /bin/bash \
  direvius/yandex-tank