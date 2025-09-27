locals {
  base_tags = {
    Environment = var.environment
    CostCenter  = "Development"
    Team        = "DevOps"
    ManagedBy   = "Terraform"
  }
  tags = merge(local.base_tags, var.tags)
}
