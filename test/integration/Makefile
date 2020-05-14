.DEFAULT_GOAL				:= all

# This target (taken from: https://gist.github.com/prwhite/8168133) is an easy way to print out a usage/ help of all make targets.
# For all make targets the text after \#\# will be printed.
help: ## Prints the help
	@echo "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\1\:\2/' | column -c2 -t -s :)"

build.0.9.1: ## builds and tags nomad in version 0.9.1 with consul in version 1.5.0
	@echo "--> build image with nomad 0.9.1 and consul 1.5.0"
	@docker build -t thobe/nomad -f Dockerfile --build-arg NOMAD_VERSION=0.9.1 --build-arg CONSUL_VERSION=1.5.0 .
	@echo "--> tag image as thobe/nomad:0.9.1"
	@docker tag thobe/nomad:latest thobe/nomad:0.9.1

push.0.9.1: build.0.9.1 ## pushes the nomad in version 0.9.1 with consul in version 1.5.0 to docker hub
	@echo "--> push image thobe/nomad:0.9.1"
	@docker push thobe/nomad:0.9.1

run.0.9.1: build.0.9.1 ## run nomad in version 0.9.1 with consul in version 1.5.0
	@echo "--> run nomad-consul container (nomad 0.9.1, consul 1.5.0)"
	@docker run --privileged -it -p 8500:8500 -p 4646:4646 thobe/nomad:0.9.1
