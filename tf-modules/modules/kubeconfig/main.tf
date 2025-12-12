resource "null_resource" "create_kubeconfig" {
  triggers = {
    cluster_id = var.cluster_id
  }

  provisioner "local-exec" {
    command = "yc managed-kubernetes cluster get-credentials --id ${var.cluster_id} --external --force --kubeconfig ${path.module}/kubeconfig.yaml"
  }
}

output "kubeconfig_path" {
  value = "${path.module}/kubeconfig.yaml"
}
