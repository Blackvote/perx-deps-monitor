resource "yandex_vpc_network" "k8s" {
  name = "perx-k8s-network"

  labels = {
    env     = "prod"
    project = "perx"
    role    = "network"
  }
}

resource "yandex_vpc_subnet" "public_a" {
  name           = "perx-k8s-public-a"
  zone           = "ru-central1-a"
  network_id     = yandex_vpc_network.k8s.id
  v4_cidr_blocks = ["10.0.0.0/24"]

  labels = {
    env     = "prod"
    project = "perx"
    type    = "subnet"
    scope   = "public"
  }
}

resource "yandex_vpc_gateway" "nat" {
  name = "perx-nat-gw"
  shared_egress_gateway {}

  labels = {
    env     = "prod"
    project = "perx"
    role    = "nat"
    scope   = "private"
  }
}

resource "yandex_vpc_route_table" "private_rt" {
  name       = "perx-private-rt"
  network_id = yandex_vpc_network.k8s.id

  static_route {
    destination_prefix = "0.0.0.0/0"
    gateway_id         = yandex_vpc_gateway.nat.id
  }

  labels = {
    env     = "prod"
    project = "perx"
    scope   = "nat"
    role    = "route-table"
  }
}

resource "yandex_vpc_subnet" "private_a" {
  name           = "perx-k8s-private-a"
  zone           = "ru-central1-a"
  network_id     = yandex_vpc_network.k8s.id
  v4_cidr_blocks = ["10.0.1.0/24"]
  route_table_id = yandex_vpc_route_table.private_rt.id

  labels = {
    env     = "prod"
    project = "perx"
    type    = "subnet"
    scope   = "private"
  }
}
