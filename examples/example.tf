# RackN 2023
# Digital Rebar v4.4+ Terraform v0.13+ Provider

terraform {
  required_version = ">= 0.13.0"
  required_providers {
    drp = {
      version = "2.3.1"
      source  = "rackn/drp"
    }
  }
}

provider "drp" {
  username = "rocketskates"
  password = "r0cketsk8ts"
  endpoint = "https://192.168.1.93:8092"
  # token  = will read from RS_TOKEN if set
  # key    = will read from RS_KEY if set
}

resource "drp_machine" "one_random_node" {

  # Required values
  # there are none!

  # Settable values
  # pool = name of an existing DRP pool (defaults to "default")
  # allocate_workflow = Name of workflow to set when a machine is allocated from the pool
  # if none is set the default pool/workflow defined by the drp admin will be used
  # deallocate_workflow = Name of workflow to set when the machine is released back to the pool
  # setting this overrides the defaults defined by the DRP admin
  # timeout = time string for max wait time (default to 5m)
  #
  # List of public SSH keys to be installed (written as Param.access-keys)
  # authorized_keys = ["ssh key"]
  #
  # List of profiles to apply to node (must already exist)
  # add_profiles = ["mandy", "clause"]
  #
  # list of parameters to set with their string value forms
  # add_parameters = ["param1: value1", "param2: value2"]
  #
  # Use filters to hone your machine to a specific set of criteria, or exclude them based on criteria
  # follows the Digital Rebar CLI command line pattern
  # This example looks for a machine named esxi-7-testing.example.local where the "Address" field on the
  # DRP Machine object is not empty/blank
  # filters = ["Name=esxi-7-testing.example.local", "Address=Ne()"]
  # This example excludes machines that do not have an Address
  # filters = ["Address=Ne()"]
  # Returned values
  # name = machine name
  # address = machine address
  # status = machine status (typically "InUse")
  # If your machine does not have an Address it will report the "address" as nil/null
}

output "machine_ip" {
  value       = drp_machine.one_random_node.address
  description = "Machine.Address (the Machine's primary IP)"
}

output "machine_id" {
  value       = drp_machine.one_random_node.id
  description = "Machine.Uuid"
}

output "machine_name" {
  value       = drp_machine.one_random_node.name
  description = "Machine.Name"
}

output "machine_status" {
  value       = drp_machine.one_random_node.status
  description = "Machine.PoolStatus"
}
