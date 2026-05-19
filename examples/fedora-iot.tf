# Butane Config for Fedora IoT
data "ct_config" "fedora-iot-config" {
  content = templatefile("${path.module}/content/fiot.yaml", {
    message = "Hello World!"
  })
  strict             = true
  pretty_print       = true
  use_mapped_version = true

  snippets = [
    file("${path.module}/content/fiot-snippet.yaml"),
  ]
}

# Render as Ignition
resource "local_file" "fedora-iot" {
  content  = data.ct_config.fedora-iot-config.rendered
  filename = "${path.module}/output/fedora-iot.ign"
}
