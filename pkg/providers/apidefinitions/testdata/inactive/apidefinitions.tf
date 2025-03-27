module "definition" {
  source      = "./modules/definition"
  contract_id = var.contract_id
  group_id    = var.group_id
}

#module "activation_staging" {
#    source      = "./modules/activation"
#    depends_on  = [module.definition]
#    api_id      = module.definition.api_id
#    api_version = module.definition.api_latest_version
#    network     = "STAGING"
#}

#module "activation_production" {
#    source      = "./modules/activation"
#    depends_on  = [module.definition]
#    api_id      = module.definition.api_id
#    api_version = module.definition.api_latest_version
#    network     = "PRODUCTION"
#}

#module "operations" {
#    source              = "./modules/definition"
#}