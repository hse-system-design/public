# Kubernetes Deployments

### Rolling update

С помощью deployment можно управлять тем, как именно мы заменяем старый реплика-сет на новый

Поднимим кубики
```bash
minikube start --vm-driver=virtualbox
minikube dashboard --url=true
```

Просто накатим версию сервиса. На данный момент никакой особой разницы с просто созданием реплика-сета нет
```bash
kubectl apply -f echo-deployment-1.yaml --record
kubectl apply -f echo-svc.yaml
curl http://`minikube ip`:30030
```

В отдельном терминале будем следить за ситуацией в кластере
```bash
kubectl get pods --show-labels -w
```

Накатываем новый деплоймент с простым способом обновления - Recreate
```bash
kubectl apply -f echo-deployment-2.yaml --record
```

Пробуем более плавный передеплой - RollingUpdate
```bash
kubectl apply -f echo-deployment-3.yaml --record
```

```bash
kubectl rollout history deployment.v1.apps/echo-service-deployment

kubectl rollout undo deployment.v1.apps/echo-service-deployment
kubectl rollout undo deployment.v1.apps/echo-service-deployment --to-revision=2
```

```bash
kubectl apply -f echo-probes-1.yaml
```