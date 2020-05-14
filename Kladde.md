./install-nomad --version 0.9.1

nomad agent -config=nomad.hcl

consul agent -config-file=consul.hcl

docker run -it -p 8500:8500 -p 4646:4646 nomad:latest

ip -o -4 addr list eth0 | head -n1 | awk '{print \$4}' | cut -d/ -f1
