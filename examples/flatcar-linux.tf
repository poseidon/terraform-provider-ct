# Butane Config for Flatcar Linux
data "ct_config" "flatcar-config" {
  content = templatefile("${path.module}/content/flatcar.yaml", {
    message = "Hello World!"
  })
  strict       = true
  pretty_print = true

  snippets = [
    file("${path.module}/content/flatcar-snippet.yaml"),
  ]
}

# Render as Ignition
resource "local_file" "flatcar" {
  content  = data.ct_config.flatcar-config.rendered
  filename = "${path.module}/output/flatcar.ign"
}
