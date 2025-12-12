resource "yandex_iam_service_account" "k8s_cluster_sa" {
  name        = "perx-k8s-cluster-sa"
  description = "Service account for Managed K8s master nodes"
}

resource "yandex_iam_service_account" "k8s_nodes_sa" {
  name        = "perx-k8s-nodes-sa"
  description = "Service account for Managed K8s worker nodes"
}

resource "yandex_resourcemanager_folder_iam_binding" "k8s_cluster_agent" {
  folder_id = var.folder_id
  role      = "k8s.clusters.agent"

  members = [
    "serviceAccount:${yandex_iam_service_account.k8s_cluster_sa.id}",
  ]
}

resource "yandex_resourcemanager_folder_iam_binding" "k8s_vpc_public_admin" {
  folder_id = var.folder_id
  role      = "vpc.publicAdmin"

  members = [
    "serviceAccount:${yandex_iam_service_account.k8s_cluster_sa.id}",
  ]
}

resource "yandex_resourcemanager_folder_iam_binding" "k8s_nodes_cr_puller" {
  folder_id = var.folder_id
  role      = "container-registry.images.puller"

  members = [
    "serviceAccount:${yandex_iam_service_account.k8s_nodes_sa.id}",
  ]
}
