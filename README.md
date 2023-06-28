# terraform-provider-rotating
Terraform Provider for rotating blue/green resources.

This provider is good for when you need to keep an active and inactive version of a resource, 
and you want to be able to rotate between them while keeping traffic/secrets/anything healthy while you rotate the inactive set.

See an example in the [examples](examples/) directory. A good rundown of how to use this provider is in the [demo](examples/demo/README.md) directory.
