# Fedora CoreOS Config -> Ignition
resource "local_file" "fedora-coreos-ign" {
  content = data.ct_config.fedora-coreos-config.rendered
  filename = "${path.module}/output/fedora-coreos.ign"
}

# Fedora CoreOS Config
data "ct_config" "fedora-coreos-config" {
  content      = data.template_file.fedora-coreos-worker.rendered
  pretty_print = true

  snippets = [
    file("${path.module}/content/fcc-snippet.yaml"),
  ]
}

# Content (e.g. possibly templated)
data "template_file" "fedora-coreos-worker" {
  template = file("${path.module}/content/fcc.yaml")

  vars = {
    message = "Hello World!"
  }
}
