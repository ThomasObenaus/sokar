#!/usr/bin/dumb-init /bin/bash
set -o errexit

consul agent -config-file=consul.hcl &
nomad agent -config=nomad.hcl