module "module_a" {
  source = "git@github.com:dhth/infrastructure//modules/applications/module-a?ref=module-a-v1.0.22"
  environment                  = var.environment
  prefix                       = var.prefix
}

module "module_b" {
  source = "git@github.com:dhth/infrastructure//modules/applications/module-b?ref=module-b-v0.1.8"
  environment                  = var.environment
  prefix                       = var.prefix
}

module "module_c" {
  source = "git@github.com:dhth/infrastructure//modules/applications/module-c?ref=module-c-v0.1.0"
  environment                  = var.environment
  prefix                       = var.prefix
}

module "module_d" {
  source = "git@github.com:dhth/infrastructure//modules/applications/module-d?ref=module-c-v0.2.0"
  environment                  = var.environment
  prefix                       = var.prefix
}
