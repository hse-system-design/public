# Kubernetes Intro

Слайды презентации можно найти в папочке [slides](../slides).

Для локальной работы с Kubernetes необходимо `minikube` и `kubectl`. 
Инструкции по установке можно найти [тут](https://kubernetes.io/docs/tasks/tools/).

## Запуск контейнеров в Kubernetes

Запускаем локальный Kubernetes в виртуалке

```
minikube delete # если до этого игрались с minikube
minikube start --vm-driver=virtualbox # если не работает на маке, попробуйте hyperkit
```

- запускаем под с echo сервером

```
kubectl run echo-server --image=ealen/echo-server
kubectl get pods -o yaml # подробное описание пода, в конце есть его IP
```


- запускаем еще один под с убунтой

```
kubectl run -it ubuntu --image=ubuntu:20.04
# внутри пода
apt update && apt install curl
curl 127.0.0.3/hello # здесь 127.0.0.3 --- IP пода, которое мы получили на предыдущем шаге
```

## Добавление репликасета

- Добавляем репликасет

```
kubectl create -f echo-rs.yaml
```

- изменяем количество реплик пода

```
kubectl scale --replicas=20 rs echo-service
```

## Добавление сервиса

- добавляем сервис, ссылающиеся на эхо-поды

```
kubectl create -f echo-svc.yaml
```

- проверяем доступность подов по DNS имени с пода ubuntu 

```
kubectl exec -it ubuntu -- bash
# внутри пода ubuntu
apt update && apt install -y curl
curl echo-service:8080/foo
```

## NodePort для доступа извне кластера

```
kubectl create -f echo-nodeport.yaml
curl `minikube ip`:30030/pupa/lupa
```

## Полезные ссылки

- интерактивный туториал [Learn Kubernetes Basics](https://kubernetes.io/docs/tutorials/kubernetes-basics/)