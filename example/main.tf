data "template_file" "fuze" {
  template = "${file("${path.module}/fuze.tpl")}"

  vars {
    message = "Hello World!"
  }
}

data "fuze_config" "example" {
  pretty_print = true
  content      = "${data.template_file.fuze.rendered}"
}

resource "null_resource" "echo" {
  triggers {
    fuze_config = "${data.fuze_config.example.id}"
  }

  provisioner "local-exec" {
    command = "echo '${data.fuze_config.example.rendered}'"
  }
}
