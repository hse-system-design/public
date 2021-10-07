Соберем образ с приложением шифратором

```bash
docker build -t crypt:latest .
```

Посмотреть на образы, которые есть сейчас на машине

```bash
docker images
```

Запустим приложение шифратор. Он шифрует файлы, поэтому для того, чтобы подать этому приложению на вход данные,
пробросим файл из основной операционной системы по адресу "./Dockerfile" внутрь контейнера в место "/data/file".

```bash
docker run -v `pwd`/Dockerfile:/data/file crypt:latest encrypt /data/file

```

Попробуем расшифровать данные. Сохраним результат работы в отдельный файл, пробросим его внутрь и дешифруем.

```bash
docker run -v `pwd`/Dockerfile:/data/file crypt:latest encrypt /data/file > enc.txt
docker run -v `pwd`/enc.txt:/biba/kuka.txt crypt:latest decrypt /biba/kuka.txt
```

Вуаля!