resource "tls_private_key" "this" {
  count     = var.use_ssh_key ? 0 : 1
  algorithm = "RSA"
  rsa_bits  = 4096
}
