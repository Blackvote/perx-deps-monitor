provider "yandex" {
  cloud_id  = var.cloud_id
  folder_id = var.folder_id
  zone      = "ru-central1-a"
}

module "network" {
  source = "./modules/network"
}

module "iam" {
  source    = "./modules/iam"
  folder_id = var.folder_id
}

module "k8s" {
  source = "./modules/k8s"

  folder_id       = var.folder_id
  network_id      = module.network.network_id
  public_subnet_id  = module.network.public_subnet_id
  public_subnet_zone = module.network.public_subnet_zone
  private_subnet_id = module.network.private_subnet_id
  private_subnet_zone = module.network.private_subnet_zone

  k8s_cluster_sa_id = module.iam.k8s_cluster_sa_id
  k8s_nodes_sa_id   = module.iam.k8s_nodes_sa_id

  k8s_api_cidrs = var.k8s_api_cidrs
}

module "kubeconfig" {
  source     = "./modules/kubeconfig"
  cluster_id = module.k8s.cluster_id

  depends_on = [module.k8s]
}

output "kubeconfig_path" {
  value = module.kubeconfig.kubeconfig_path
}
