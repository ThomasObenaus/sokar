#!/usr/bin/dumb-init /bin/bash
set -o errexit

NOMAD_VERSION=$(cat nomad.version)
CONSUL_VERSION=$(cat consul.version)


function replace_template_var_in_files {
  local readonly workingDir=$1
  local readonly templateToReplace=$2
  local readonly value=$3

  files=(consul.hcl nomad.hcl)

  for item in ${files[*]}
  do
    file=${workingDir}/${item}
    echo "Replacing ${templateToReplace} by ${value} in ${file}."
    command="sed -i.bak s/${templateToReplace}/${value}/g ${file}"
    eval ${command}
  done
}

# patches the nomad.hcl and consul.hcl to use the local docker ip address
# usually its 172.17.0.2 but not in each case since it depends on the local docker configuration.
function patch_config_files {
    echo "[PATCHING CONFIG FILES] (with correct ip address) ######################"

    workingDir=$(pwd)
    ipAddr=$(ip -o -4 addr list eth0 | head -n1 | awk '{print $4}' | cut -d/ -f1)

    echo "[PATCHING CONFIG FILES] ipAddr=$ipAddr"
    replace_template_var_in_files "${workingDir}" "{{host_ip_address}}" "${ipAddr}"

    echo "[PATCHING CONFIG FILES] patched files"
    cat nomad.hcl
    cat consul.hcl

    echo "[PATCHING CONFIG FILES] (with correct ip address) ######################"
}

patch_config_files
sleep 1


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