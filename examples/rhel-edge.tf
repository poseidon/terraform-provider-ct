# Butane Config for RHEL for Edge
data "ct_config" "rhel-edge-config" {
  content = templatefile("${path.module}/content/r4e.yaml", {
    message = "Hello World!"
  })
  strict       = true
  pretty_print = true

  snippets = [
    file("${path.module}/content/r4e-snippet.yaml"),
  ]
}

# Render as Ignition
resource "local_file" "rhel-edge" {
  content  = data.ct_config.rhel-edge-config.rendered
  filename = "${path.module}/output/rhel-edge.ign"
}
