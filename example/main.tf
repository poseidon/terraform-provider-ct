data "template_file" "ct" {
  template = "${file("${path.module}/ct.tpl")}"

  vars {
    message = "Hello World!"
  }
}

data "ct_config" "example" {
  pretty_print = true
  content      = "${data.template_file.ct.rendered}"
}

resource "null_resource" "echo" {
  triggers {
    ct_config = "${data.ct_config.example.id}"
  }

  provisioner "local-exec" {
    command = "echo '${data.ct_config.example.rendered}'"
  }
}
