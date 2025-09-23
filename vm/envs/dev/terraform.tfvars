prefix      = "taller1"
location    = "eastus"
environment = "dev"

vm_size        = "Standard_B1s"
admin_username = "ubuntu"

use_ssh_key         = false
ssh_public_key_path = null

allowed_ssh_cidrs = ["0.0.0.0/0"]

tags = {
  CostCenter = "Development"
  Team       = "DevOps"
}
