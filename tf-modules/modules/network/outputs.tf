output "network_id" {
  value = yandex_vpc_network.k8s.id
}

output "public_subnet_id" {
  value = yandex_vpc_subnet.public_a.id
}

output "public_subnet_zone" {
  value = yandex_vpc_subnet.public_a.zone
}

output "private_subnet_id" {
  value = yandex_vpc_subnet.private_a.id
}

output "private_subnet_zone" {
  value = yandex_vpc_subnet.private_a.zone
}
