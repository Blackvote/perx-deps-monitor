output "k8s_cluster_sa_id" {
  value = yandex_iam_service_account.k8s_cluster_sa.id
}

output "k8s_nodes_sa_id" {
  value = yandex_iam_service_account.k8s_nodes_sa.id
}
