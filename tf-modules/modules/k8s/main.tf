resource "yandex_vpc_security_group" "k8s_cluster_api" {
  name       = "perx-k8s-cluster-api"
  network_id = var.network_id

  ingress {
    description    = "K8s API via 443"
    protocol       = "TCP"
    port           = 443
    v4_cidr_blocks = var.k8s_api_cidrs
  }

  ingress {
    description    = "K8s API via 6443"
    protocol       = "TCP"
    port           = 6443
    v4_cidr_blocks = var.k8s_api_cidrs
  }

  egress {
    description    = "Allow egress (control plane needs it)"
    protocol       = "ANY"
    from_port      = 0
    to_port        = 65535
    v4_cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "yandex_vpc_security_group" "k8s_cluster_node_traffic" {
  name       = "perx-k8s-cluster-node-traffic"
  network_id = var.network_id

  ingress {
    description       = "NLB health checks"
    protocol          = "TCP"
    from_port         = 0
    to_port           = 65535
    predefined_target = "loadbalancer_healthchecks"
  }

  ingress {
    description       = "Master<->Nodes service traffic (ingress)"
    protocol          = "ANY"
    from_port         = 0
    to_port           = 65535
    predefined_target = "self_security_group"
  }

  egress {
    description       = "Master<->Nodes service traffic (egress)"
    protocol          = "ANY"
    from_port         = 0
    to_port           = 65535
    predefined_target = "self_security_group"
  }
}

resource "yandex_vpc_security_group" "k8s_nodes_ingress" {
  name       = "perx-k8s-nodes-ingress"
  network_id = var.network_id

  ingress {
    description    = "Ingress via NodePort from LB / Internet"
    protocol       = "TCP"
    from_port      = 30000
    to_port        = 32767
    v4_cidr_blocks = var.k8s_api_cidrs
  }
}

resource "yandex_vpc_security_group" "k8s_nodes_base" {
  name       = "perx-k8s-nodes-base"
  network_id = var.network_id

  ingress {
    description    = "Pods<->Services (cluster/service CIDRs)"
    protocol       = "ANY"
    from_port      = 0
    to_port        = 65535
    v4_cidr_blocks = ["10.10.0.0/16", "10.11.0.0/16"]
  }

  egress {
    description    = "Nodes egress (images, updates, etc)"
    protocol       = "ANY"
    from_port      = 0
    to_port        = 65535
    v4_cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "yandex_kubernetes_cluster" "perx" {
  name       = "perx-cluster"
  folder_id  = var.folder_id
  network_id = var.network_id

  cluster_ipv4_range = "10.10.0.0/16"
  service_ipv4_range = "10.11.0.0/16"

  service_account_id      = var.k8s_cluster_sa_id
  node_service_account_id = var.k8s_nodes_sa_id

  release_channel         = "REGULAR"
  network_policy_provider = "CALICO"

  node_ipv4_cidr_mask_size = 26

  master {
    version   = "1.33"
    public_ip = true

    zonal {
      zone      = var.public_subnet_zone
      subnet_id = var.public_subnet_id
    }

    security_group_ids = [
    yandex_vpc_security_group.k8s_cluster_api.id,
    yandex_vpc_security_group.k8s_cluster_node_traffic.id,
  ]

    maintenance_policy {
      auto_upgrade = false
    }
  }

  labels = {
    env     = "prod"
    project = "perx"
    role    = "k8s-cluster"
    scope   = "public"
  }
}

resource "yandex_kubernetes_node_group" "perx_ng" {
  cluster_id = yandex_kubernetes_cluster.perx.id

  name        = "perx-node-group-1"
  description = "Node group"
  version     = "1.33"

  labels = {
    env     = "prod"
    project = "perx"
    role    = "worker"
  }

  scale_policy {
    fixed_scale {
      size = 1
    }
  }

  allocation_policy {
    location {
      zone = var.private_subnet_zone
    }
  }

  maintenance_policy {
    auto_upgrade = false
    auto_repair  = false
  }

  instance_template {
    platform_id = "standard-v3"

    resources {
      cores  = 2
      memory = 4
    }

    boot_disk {
      type = "network-hdd"
      size = 50
    }

    network_interface {
      subnet_ids         = [var.private_subnet_id]
      nat                = false
      security_group_ids = [
      yandex_vpc_security_group.k8s_cluster_node_traffic.id,
      yandex_vpc_security_group.k8s_nodes_base.id,
      yandex_vpc_security_group.k8s_nodes_ingress.id,
    ]
    }
  }
}
