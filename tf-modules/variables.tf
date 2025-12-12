variable "cloud_id" {
  type        = string
  description = "Yandex Cloud ID"
}

variable "folder_id" {
  type        = string
  description = "Yandex Folder ID"
}

variable "k8s_api_cidrs" {
  description = "CIDR blocks allowed to access NodePort / ingress"
  type        = list(string)
  default     = ["0.0.0.0/0"]
}