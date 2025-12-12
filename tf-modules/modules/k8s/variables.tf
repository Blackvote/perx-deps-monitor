variable "folder_id" {
  type        = string
  description = "Yandex Folder ID"
}

variable "network_id" {
  type        = string
  description = "VPC network id"
}

variable "public_subnet_id" {
  type        = string
  description = "Public subnet id for master"
}

variable "public_subnet_zone" {
  type        = string
  description = "Public subnet zone for master"
}

variable "private_subnet_id" {
  type        = string
  description = "Private subnet id for worker nodes"
}

variable "private_subnet_zone" {
  type        = string
  description = "Private subnet zone for worker nodes"
}

variable "k8s_cluster_sa_id" {
  type        = string
  description = "Service account id for cluster"
}

variable "k8s_nodes_sa_id" {
  type        = string
  description = "Service account id for nodes"
}

variable "k8s_api_cidrs" {
  description = "CIDR blocks allowed to access NodePort / ingress"
  type        = list(string)
  default     = ["0.0.0.0/0"]
}