variable "PRIVX_API_BASE_URL" {
  type        = string
  default     = "https://privx.pf-beta.cloudservicesfactory.com"
  description = "privx api base url"
}

variable "PRIVX_OAUTH_CLIENT_ID" {
  type        = string
  description = "privx api oauth client ID"
}

variable "PRIVX_OAUTH_CLIENT_SECRET" {
  type        = string
  description = "privx api oauth client secret"
}

variable "PRIVX_API_CLIENT_ID" {
  type        = string
  description = "privx api client id"
}

variable "PRIVX_API_CLIENT_SECRET" {
  type        = string
  description = "privx api client id"
}

variable "PRIVX_DEBUG" {
  type        = bool
  default     = false
  description = "Privx provider debug mode"
}
