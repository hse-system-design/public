```bash
docker build -t crypt:latest .
```

```bash
docker run -v `pwd`/Dockerfile:/data/file crypt:latest encrypt /data/file

docker run -v `pwd`/Dockerfile:/data/file crypt:latest encrypt /data/file > enc.txt
docker run -v `pwd`/enc.txt:/biba/kuka.txt crypt:latest decrypt /biba/kuka.txt
```