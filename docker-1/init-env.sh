#!/bin/bash

set -xe

mkdir -p crypt

curl https://raw.githubusercontent.com/ADKosm/system-design-2021-public/main/docker-1/crypt/crypt.py > crypt/crypt.py
curl https://raw.githubusercontent.com/ADKosm/system-design-2021-public/main/docker-1/crypt/Dockerfile > crypt/Dockerfile

mkdir -p server

curl https://raw.githubusercontent.com/ADKosm/system-design-2021-public/main/docker-1/server/server.py > server/server.py
curl https://raw.githubusercontent.com/ADKosm/system-design-2021-public/main/docker-1/server/Dockerfile > server/Dockerfile

mkdir -p compose

curl https://raw.githubusercontent.com/ADKosm/system-design-2021-public/main/docker-1/compose/docker-compose.yaml > compose/docker-compose.yaml
curl https://raw.githubusercontent.com/ADKosm/system-design-2021-public/main/docker-1/compose/Dockerfile.server > compose/Dockerfile.server
curl https://raw.githubusercontent.com/ADKosm/system-design-2021-public/main/docker-1/compose/Dockerfile.vault > compose/Dockerfile.vault
curl https://raw.githubusercontent.com/ADKosm/system-design-2021-public/main/docker-1/compose/secret.key > compose/secret.key
curl https://raw.githubusercontent.com/ADKosm/system-design-2021-public/main/docker-1/compose/server.py > compose/server.py
curl https://raw.githubusercontent.com/ADKosm/system-design-2021-public/main/docker-1/compose/vault.py > compose/vault.py
