variable "prefix" { type = string }
variable "location" { type = string }
variable "rg_name" { type = string }
variable "allowed_ssh" { type = list(string) }
variable "tags" { type = map(string) }
