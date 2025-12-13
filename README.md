## Подключение к Yandex Cloud

### 1. Авторизация yc

```bash
yc init
```

Проверь:

```bash
yc config list
```

Должны быть заданы:

* `cloud-id`
* `folder-id`
* `compute-default-zone`

Синтаксис: `yc config set <key> <value>`

---

### 2. Переменные окружения для Terraform

Создай файл `.env` (или экспортируй переменные):

```bash
export YC_TOKEN=$(yc iam create-token --impersonate-service-account-id <SA_ID>)
export YC_CLOUD_ID=<your_cloud_id>
export YC_FOLDER_ID=<your_folder_id>
```
---

## Развёртывание инфраструктуры (Terraform)

Перейди в каталог Terraform:

```bash
cd tf-modules
terraform init
terraform plan
terraform apply
```

В результате будут созданы:

* VPC и подсети
* NAT Gateway
* Security Groups
* Managed Kubernetes Cluster
* Node Group


## Установка ingress-nginx

Добавь Helm-репозиторий:

```bash
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm install ingress-nginx ingress-nginx/ingress-nginx \
  --namespace ingress-nginx \
  --create-namespace
```

---

## Создание Kubernetes Secrets

Перед установкой приложения необходимо создать секреты для MongoDB и NATS:

```bash
kubectl create secret generic perx-mongodb-creds \
  --from-literal=mongodb-passwords='perx-password'

kubectl create secret generic perx-nats-creds \
  --from-literal=password='perx-password'
```

## Деплой приложения (Helm)

Перейди в каталог Helm-чарта:

```bash
cd helm
helm dependency update
helm install perx . \
  --namespace <NS>
```
