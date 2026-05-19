# Butane Config for OpenShift
data "ct_config" "openshift-config" {
  content = templatefile("${path.module}/content/openshift.yaml", {
    message = "Hello World!"
  })
  strict       = true
  pretty_print = true
}

# Render as MachineConfig
resource "local_file" "openshift" {
  content  = data.ct_config.openshift-config.rendered
  filename = "${path.module}/output/openshift.yaml"
}
