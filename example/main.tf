# Container Linux config
data "ct_config" "worker" {
  pretty_print = true
  content      = "${data.template_file.worker.rendered}"

  snippets = [
    "${file("${path.module}/snippet.yaml")}",
  ]
}

# Content (e.g. templated)
data "template_file" "worker" {
  template = "${file("${path.module}/clc.tmpl")}"

  vars = {
    message = "Hello World!"
  }
}

# Example usage
resource "null_resource" "echo" {
  triggers = {
    ct_config = "${data.ct_config.worker.id}"
  }

  provisioner "local-exec" {
    command = "echo '${data.ct_config.worker.rendered}'"
  }
}
