output "public_ip" {
  description = "IP pública de la VM"
  value       = module.compute.public_ip
}

# útil para copiar al Secret PROD_SSH_KEY cuando use_ssh_key=false
output "private_key_pem" {
  value     = module.keypair.private_key_pem
  sensitive = true
}
