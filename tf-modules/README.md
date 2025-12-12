# perx-terraform-modular

Modules:

- `modules/network` — VPC + public/private subnets + NAT gateway + route table
- `modules/iam` — service accounts + folder IAM bindings
- `modules/k8s` — security group + managed k8s cluster + node group
- `modules/kubeconfig` — генерация kubeconfig через `yc managed-kubernetes ...` 
