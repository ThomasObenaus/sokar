
DASHBOARD_DIR="../../dashboards"

help:
	@echo "Available make targets:"
	@echo "\t- up\t\t\tStarts up a prometheus and a grafana instance,"
	@echo "\t\t\t\tscraping metrics of sokar and providing a dashboard for sokar."
	@echo "\t- down\t\t\tStops the monitoring components (prometheus, grafana)."

up: sep
	@echo "--> Startup (+build) monitoring components"
	@echo "--> Copy dashboards from $(DASHBOARD_DIR)"
	@mkdir -p grafana/dashboards
	@cp -R $(DASHBOARD_DIR) grafana
	@docker-compose up --build -d
	@rm -rf grafana/dashboards
	@xdg-open http://localhost:3000

down: sep
	@echo "--> Stop monitoring components"
	@docker-compose down

sep:
	@echo "----------------------------------------------------------------------------------"