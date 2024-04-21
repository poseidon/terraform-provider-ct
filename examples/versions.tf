
terraform {
  required_version = ">= 0.13.0"
  required_providers {
    local = "~> 2.0"
    ct = {
      source  = "poseidon/ct"
      version = "~> 0.13.0"
      #source  = "terraform.localhost/poseidon/ct"
      #version = "0.12.0"
    }
  }
}


