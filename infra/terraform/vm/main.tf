module "rg" {
  source   = "./modules/resource_group"
  prefix   = var.prefix
  location = var.location
  tags     = local.tags
}

module "keypair" {
  source              = "./modules/keypair"
  use_ssh_key         = var.use_ssh_key
  ssh_public_key_path = var.ssh_public_key_path
}

module "network" {
  source      = "./modules/network"
  prefix      = var.prefix
  location    = var.location
  rg_name     = module.rg.name
  allowed_ssh = var.allowed_ssh_cidrs
  tags        = local.tags
}

module "compute" {
  source         = "./modules/compute_linux_vm"
  prefix         = var.prefix
  location       = var.location
  rg_name        = module.rg.name
  subnet_id      = module.network.subnet_id
  nsg_id         = module.network.nsg_id
  admin_username = var.admin_username
  vm_size        = var.vm_size
  # pública final que usa la VM
  public_key = var.use_ssh_key ? file(var.ssh_public_key_path) : module.keypair.public_key_openssh
  tags       = local.tags
}
