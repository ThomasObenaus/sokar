
DASHBOARD_DIR="../../dashboards"

# This target (taken from: https://gist.github.com/prwhite/8168133) is an easy way to print out a usage/ help of all make targets.
# For all make targets the text after \#\# will be printed.
help: ## Prints the help
	@echo "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\1\:\2/' | column -c2 -t -s :)"


up: sep ## Start prometheus and grafana locally
	@echo "--> Startup (+build) monitoring components"
	@echo "--> Copy dashboards from $(DASHBOARD_DIR)"
	@mkdir -p grafana/dashboards
	@cp -R $(DASHBOARD_DIR) grafana
	@export HOST_IP_NOT_SET=$$(ifconfig | grep docker0 | wc -l); docker-compose up --build -d
	@rm -rf grafana/dashboards
	@xdg-open http://localhost:3000

down: sep ## Stop prometheus and grafana
	@echo "--> Stop monitoring components"
	@export HOST_IP_NOT_SET=0; docker-compose down --remove-orph

sep:
	@echo "----------------------------------------------------------------------------------"