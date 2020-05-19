#!/bin/sh

# Check if host.docker.internal is already set (e.g. on windows it is prepared already)
if [[ $HOST_IP_NOT_SET -eq 1 ]]
then  
  echo "Settig up host.docker.internal since it is not yet set"
  echo "Adding host.docker.internal to /etc/hosts"
  ip -4 route list match 0/0 | awk '{print $3 " host.docker.internal"}' >> /etc/hosts
else
  echo "No need to set host.docker.internal, its already set."
fi

# redo the check
echo "Resolved address for host.docker.internal:"
ping -c 1 host.docker.internal

# call original (delegate back)
echo "Calling original with $@"
/bin/prometheus "$@"