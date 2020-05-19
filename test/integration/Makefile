.DEFAULT_GOAL				:= test.complete-setup

# This target (taken from: https://gist.github.com/prwhite/8168133) is an easy way to print out a usage/ help of all make targets.
# For all make targets the text after \#\# will be printed.
help: ## Prints the help
	@echo "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\1\:\2/' | column -c2 -t -s :)"

test.complete-setup: ## Runs the integration-test and ensures that sokar is running as well
	@echo "--> Starting integration-tests of sokar and thobe/sokar:latest as test-object"
	@docker-compose up --build --abort-on-container-exit --exit-code-from integration_test
	@docker-compose down

########################################### TestCases ##############################################################
tc.integration: ## Run integration-test test-cases endpoints
	@echo "############# Run Integration Tests ........ ####################################################"
	SOKAR_ADDR=http://sokar:11000 NOMAD_ADDR=http://nomad:4646 go test -v -count=1 .
	@echo "############# Run Integration Tests ........ ####################################################"
