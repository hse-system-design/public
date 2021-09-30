Чтобы запустить стенд

```bash
docker-compose up --build
```

Адреса
* localhost:3000 - grafana
* localhost:9100 - telegraf
* localhost:9090 - prometheus
* localhost:8080 - основное приложение

Чтобы запустить танк

```bash
bash run.sh
yandex-tank -c tank.yaml
```
