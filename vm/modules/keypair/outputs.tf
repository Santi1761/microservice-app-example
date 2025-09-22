output "public_key_openssh" {
  value = var.use_ssh_key ? file(var.ssh_public_key_path) : tls_private_key.this[0].public_key_openssh
}

output "private_key_pem" {
  value     = var.use_ssh_key ? "" : tls_private_key.this[0].private_key_pem
  sensitive = true
}
