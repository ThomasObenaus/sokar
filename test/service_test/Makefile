.DEFAULT_GOAL := test.complete-setup

# This target (taken from: https://gist.github.com/prwhite/8168133) is an easy way to print out a usage/ help of all make targets.
# For all make targets the text after \#\# will be printed.
help: ## Prints the help
	@echo "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\1\:\2/' | column -c2 -t -s :)"

test.complete-setup: ## Runs the service-test and ensures that sokar is running as well
	@echo "--> Starting service-tests of sokar and thobe/sokar:latest as test-object"
	@docker-compose -f endpoints.docker-compose.yml up --build --abort-on-container-exit --exit-code-from service_test
	@docker-compose -f endpoints.docker-compose.yml down

########################################### TestCases ##############################################################
tc.endpoints: ## TestCase for sokar endpoints
	@echo "############# TestCase EndPoints (_EndPoints_) ........ ####################################################"
	go test -v -run="_EndPoints_" .
	@echo "############# TestCase EndPoints (_EndPoints_) Complete ####################################################"
