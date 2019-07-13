# Container Linux Config -> Ignition
resource "local_file" "container-linux-ign" {
  content = data.ct_config.container-linux-config.rendered
  filename = "${path.module}/output/container-linux.ign"
}

# Container Linux Config
data "ct_config" "container-linux-config" {
  content      = data.template_file.container-linux-worker.rendered
  pretty_print = true

  snippets = [
    file("${path.module}/content/clc-snippet.yaml"),
  ]
}

# Content (e.g. possibly templated)
data "template_file" "container-linux-worker" {
  template = file("${path.module}/content/clc.yaml")

  vars = {
    message = "Hello World!"
  }
}
