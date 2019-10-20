.DEFAULT_GOAL				:= test.complete-setup
build_destination := "."

# This target (taken from: https://gist.github.com/prwhite/8168133) is an easy way to print out a usage/ help of all make targets.
# For all make targets the text after \#\# will be printed.
help: ## Prints the help
	@echo "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\1\:\2/' | column -c2 -t -s :)"

deps-install: ## Install the dependencies.
	@echo "--> install dependencies as listed in Gopkg.toml and Gopkg.lock"
	@dep ensure -v

test.complete-setup: ## Runs the service-test and ensures that sokar is running as well
	@echo "--> Starting service-tests of sokar and thobe/sokar:latest as test-object"
	@docker-compose -f endpoints.docker-compose.yml up --build --abort-on-container-exit --exit-code-from service_test

########################################### TestCases ##############################################################
tc.endpoints: ## TestCase for sokar endpoints
	@echo "############# TestCase EndPoints (_EndPoints_) ........ ####################################################"
	go test -v -run="_EndPoints_" .
	@echo "############# TestCase EndPoints (_EndPoints_) Complete ####################################################"


# TODO: Remove this
run:
	go test -v -count=1 -run="_My" .
