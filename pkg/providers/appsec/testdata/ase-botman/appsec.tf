module "security" {
  source      = "./modules/security"
  hostnames   = var.hostnames
  name        = var.name
  description = var.description
  contract_id = var.contract_id
  group_name  = var.group_name
}

module "activate-security" {
  source              = "./modules/activate-security"
  name                = var.name
  config_id           = module.security.config_id
  network             = var.network
  notification_emails = var.emails
  note                = var.activation_note
  depends_on          = [module.security]
}
