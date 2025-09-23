# Infraestructura como código (Terraform)

## Estructura
vm/
├─ modules/
│  ├─ resource_group/   # RG
│  ├─ network/          # VNet + Subnet + NSG
│  ├─ keypair/          # RSA 4096 si no usas clave local
│  └─ compute_linux_vm/ # Public IP + NIC + VM
├─ envs/
│  └─ dev/terraform.tfvars
├─ providers.tf / versions.tf / locals.tf / main.tf / outputs.tf / variables.tf

## Variables de ejemplo (envs/dev/terraform.tfvars)
prefix         = "taller1"
location       = "eastus"
admin_username = "ubuntu"
tags = { Environment = "dev", Team = "DevOps", ManagedBy = "Terraform", CostCenter = "Development" }

## Credenciales (en GitHub Actions o local)
ARM_CLIENT_ID, ARM_CLIENT_SECRET, ARM_TENANT_ID, ARM_SUBSCRIPTION_ID

## Comandos
terraform init
terraform plan -var-file="envs/dev/terraform.tfvars"
terraform apply -var-file="envs/dev/terraform.tfvars"

## Outputs
- public_ip
- private_key_pem (si se generó con tls_private_key)

## Conexión
ssh -i vm/private_ssh_key.pem ubuntu@$(terraform output -raw public_ip)
