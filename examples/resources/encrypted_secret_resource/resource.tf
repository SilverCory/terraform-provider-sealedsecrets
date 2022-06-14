resource "sealedsecrets_secret" "my_secret" {
  encrypted_secret = "{encrypted cool secret that lets you do cool stuff (this will be long...)}"
}

// Do something with this secret by `sealedsecrets_secret.my_secret.value`
output "cool_secret_value" {
  value = sealedsecrets_secret.my_secret.value
}
