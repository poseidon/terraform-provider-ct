
terraform {
  required_version = ">= 0.13.0"
  required_providers {
    local    = "~> 1.2"
    template = "~> 2.1"

    ct = {
      source  = "poseidon/ct"
      version = "~> 0.9.0"
      #source = "terraform.localhost/poseidon/ct"
      #version = "0.9.0"
    }
  }
}


