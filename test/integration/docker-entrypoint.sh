#!/usr/bin/dumb-init /bin/bash
set -o errexit

NOMAD_VERSION=$(cat nomad.version)
CONSUL_VERSION=$(cat consul.version)

# Taken from https://github.com/jpetazzo/dind in order to support docker in docker
echo "############################ START WRAPDOCKER #########################"
echo "#                                                                     #"
echo "#                                                                     #"
echo "#                                                                     #"
echo "#                                                                     #"
echo "#                                                                     #"
echo "#                                                                     #"
echo "#                                                                     #"
echo "############################ START WRAPDOCKER #########################"
wrapdocker &

sleep 1
echo "############################ START CONSUL ($CONSUL_VERSION) #############################"
echo "#                                                                     #"
echo "#                                                                     #"
echo "#                                                                     #"
echo "#                                                                     #"
echo "#                                                                     #"
echo "#                                                                     #"
echo "#                                                                     #"
echo "############################ START CONSUL #############################"
consul agent -config-file=consul.hcl &

sleep 1
echo "############################ START NOMAD ($NOMAD_VERSION) ######################"
echo "#                                                                     #"
echo "#                                                                     #"
echo "#                                                                     #"
echo "#                                                                     #"
echo "#                                                                     #"
echo "#                                                                     #"
echo "############################ START NOMAD ##############################"
nomad agent -config=nomad.hcl