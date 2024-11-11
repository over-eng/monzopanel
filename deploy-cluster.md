## How to provision a k8s cluster

1. First create all resources required from hetzner using `terraform apply`

2. Initialise the cluster with `pyinfra ../../../pyinfra/pyinfra_kubeadm_nodes/inventory.py pyinfra_kubeadm_nodes.deploy_initial_control_node`.
Note that the inventory will include all nodes but skip every node accept the terraform allocated initial control node.

5. Join all other nodes to the cluster `pyinfra ../../../pyinfra/pyinfra_kubeadm_nodes/inventory.py pyinfra_kubeadm_nodes.deploy_k8s_node`.
Again this will run against all nodes but the script is idempotent so it is safe to run against the initial node again.

6. Grab the admin.conf from the control plane `rsync --rsync-path="sudo rsync" ssh-user@<ip-address-control-node>:/etc/kubernetes/admin.conf .`
This is a work around until SSO with rbac is setup.

7. Apply each stage using kustomize e.g. `kustomize-build stage1 | kustomize-deploy`