variable "prefix" {
  description = "Prefijo para nombrar recursos"
  type        = string
  default     = "taller1"
}

variable "location" {
  description = "Región de Azure (ej: eastus, brazilsouth)"
  type        = string
  default     = "eastus"
}

variable "admin_username" {
  description = "Usuario administrador (debe alinear con Ansible)"
  type        = string
  default     = "ubuntu"
}

variable "vm_size" {
  description = "Tamaño de VM"
  type        = string
  default     = "Standard_B1s"
}

variable "environment" {
  description = "Ambiente (dev|stg|prod)"
  type        = string
  default     = "dev"
}

variable "allowed_ssh_cidrs" {
  description = "CIDRs que pueden hacer SSH (22). Usa 0.0.0.0/0 para todos"
  type        = list(string)
  default     = ["0.0.0.0/0"]
}

variable "use_ssh_key" {
  description = "Si true, usa ssh_public_key_path; si false, Terraform genera una RSA"
  type        = bool
  default     = false
}

variable "ssh_public_key_path" {
  description = "Ruta a la clave pública .pub si use_ssh_key=true"
  type        = string
  default     = null
}

variable "tags" {
  description = "Tags adicionales"
  type        = map(string)
  default     = {}
}
