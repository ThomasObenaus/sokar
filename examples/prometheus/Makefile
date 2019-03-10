PROMETHEUS_CONTAINER_NAME := "sokar.prometheus"

prometheus.build:
	@echo "----------------------------------------------------------------------------------"
	@echo "--> Build ${PROMETHEUS_CONTAINER_NAME} image"
	@docker build -t "${PROMETHEUS_CONTAINER_NAME}" examples/prometheus

prometheus.run:
	@echo "----------------------------------------------------------------------------------"
	@echo "--> Run ${PROMETHEUS_CONTAINER_NAME} image"
	@docker run --network="host" --name ${PROMETHEUS_CONTAINER_NAME} -d -p 127.0.0.1:9090:9090 ${PROMETHEUS_CONTAINER_NAME}

prometheus.start:
	@echo "----------------------------------------------------------------------------------"
	@echo "--> Start ${PROMETHEUS_CONTAINER_NAME} image"
	@docker start ${PROMETHEUS_CONTAINER_NAME}

prometheus.stop:
	@echo "----------------------------------------------------------------------------------"
	@echo "--> Stop ${PROMETHEUS_CONTAINER_NAME}"
	@docker stop ${PROMETHEUS_CONTAINER_NAME}