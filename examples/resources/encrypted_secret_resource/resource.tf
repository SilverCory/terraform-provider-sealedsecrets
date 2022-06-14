resource "encrypted_secret" "my_secret" {
  encrypted_secret = "{encrypted cool secret that lets you do cool stuff (this will be long...)}"
}

// Do something with this secret by `encrypted_secret.my_secret.value`
