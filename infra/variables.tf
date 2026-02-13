variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "eu-north-1"
}

variable "project_name" {
  description = "Project name for resource naming"
  type        = string
  default     = "vetles-roetter"
}

variable "turso_database_url" {
  description = "Turso database URL"
  type        = string
  sensitive   = true
}

variable "turso_auth_token" {
  description = "Turso auth token"
  type        = string
  sensitive   = true
}
