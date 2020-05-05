# Nomad Structs

Sadly due to dependency problems, especially with the replacement and patching of github.com/ugorji/go/codec to github.com/hashicorp/go-msgpack it is not possible to declare a required dependency to github.com/hashicorp/nomad in order to get the suitable nomad structs. Hence this package is just a copy of the needed code (constants, structs, types).

**The code was taken from:** https://github.com/hashicorp/nomad/blob/v0.9.1/nomad/structs/structs.go.
