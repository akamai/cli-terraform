module "definition" {
  source = "./modules/definition"
}

module "activation_staging" {
  source      = "./modules/activation"
  depends_on  = [module.definition]
  api_id      = module.definition.api_id
  api_version = module.definition.api_staging_version
  network     = "STAGING"
}

module "activation_production" {
  source      = "./modules/activation"
  depends_on  = [module.definition]
  api_id      = module.definition.api_id
  api_version = module.definition.api_production_version
  network     = "PRODUCTION"
}
