.DEFAULT_GOAL := all
name := "sokar-bin"
build_destination := "."
sokar_file_name := $(build_destination)/$(name)
docker_image := "thobe/sokar:latest"

build_time := $(shell date '+%Y-%m-%d_%H-%M-%S')
rev  := $(shell git rev-parse --short HEAD)
flag := $(shell git diff-index --quiet HEAD -- || echo "_dirty";)
tag := $(shell git describe --tags)
branch := $(shell git branch | grep \* | cut -d ' ' -f2)
revision := $(rev)$(flag)
build_info := $(build_time)_$(revision)

packages := ./scaleschedule ./config ./alertmanager ./nomad ./scaler ./helper ./scaleAlertAggregator ./sokar ./sokar/iface ./capacityplanner ./aws ./awsEc2 ./nomadWorker ./api ./

all: tools test build lint finish

# This target (taken from: https://gist.github.com/prwhite/8168133) is an easy way to print out a usage/ help of all make targets.
# For all make targets the text after \#\# will be printed.
help: ## Prints the help
	@echo "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\1\:\2/' | column -c2 -t -s :)"

.PHONY: test
test: sep gen-mocks ## Runs all unittests and generates a coverage report.
	@echo "--> Run the unit-tests"
	@go test ${packages} -race -timeout 30s -covermode=atomic -coverprofile=coverage.out

# TODO: This make target can be removed as soon as https://github.com/ThomasObenaus/sokar/issues/138 is fixed
test.no-race: sep gen-mocks ## Runs all unittests and generates a coverage report.
	@echo "--> Run the unit-tests"
	@go test ${packages} -timeout 30s -covermode=atomic -coverprofile=coverage.out

build: sep ## Builds the sokar binary.
	@echo "--> Build the $(name) in $(build_destination)"
	@go build -v -ldflags "-X main.version=$(tag) -X main.buildTime=$(build_time) -X main.revision=$(revision) -X main.branch=$(branch)" -o $(sokar_file_name) .

tools: sep ## Installs needed tools (i.e. mock generators).
	@echo "--> Install needed tools."
	@go get golang.org/x/tools/cmd/cover

gen-mocks: sep ## Generates test doubles (mocks).
	@echo "--> generate mocks (github.com/golang/mock/gomock is required for this)"
	@go get github.com/golang/mock/gomock
	@go install github.com/golang/mock/mockgen
	@mockgen -source=nomad/nomadclient_IF.go -destination test/mocks/nomad/mock_nomadclient_IF.go
	@mockgen -source=nomadWorker/nomadclient_IF.go -destination test/mocks/nomadWorker/mock_nomadclient_IF.go
	@mockgen -source=aws/iface/autoscaling_IF.go -destination test/mocks/aws/mock_autoscaling_IF.go 
	@mockgen -source=scaler/scalingtarget_IF.go -destination test/mocks/scaler/mock_scalingtarget_IF.go 
	@mockgen -source=sokar/iface/scaler_IF.go -destination test/mocks/sokar/mock_scaler_IF.go 
	@mockgen -source=sokar/iface/capacity_planner_IF.go -destination test/mocks/sokar/mock_capacity_planner_IF.go 
	@mockgen -source=sokar/iface/scaleEventEmitter_IF.go -destination test/mocks/sokar/mock_scaleEventEmitter_IF.go 
	@mockgen -source=sokar/iface/scaleschedule_IF.go -destination test/mocks/sokar/mock_scaleschedule_IF.go
	@mockgen -source=metrics/metrics.go -destination test/mocks/metrics/mock_metrics.go 
	@mockgen -source=capacityplanner/scaleschedule_IF.go -destination test/mocks/capacityplanner/mock_scaleschedule_IF.go

gen-metrics-md: sep ## Generate metrics documentation (Metrics.md) based on defined metrics in code.
	@echo "--> generate Metrics.md"
	@python3 gen_metrics_doc.py > Metrics.md

run.aws-ec2: sep build ## Builds + runs sokar locally in aws ec2 mode.
	@echo "--> Run $(sokar_file_name)"
	$(sokar_file_name) --config-file="examples/config/full.yaml" --sca.mode="aws-ec2"

run.nomad-dc: sep build ## Builds + runs sokar locally in data-center mode.
	@echo "--> Run $(sokar_file_name)"
	$(sokar_file_name) --config-file="examples/config/minimal.yaml" --sca.mode="nomad-dc"

run.nomad-job: sep build ## Builds + runs sokar locally in job mode.
	@echo "--> Run $(sokar_file_name)"
	$(sokar_file_name) --config-file="examples/config/minimal.yaml"

docker.build: sep ## Builds the sokar docker image.
	@echo "--> Build docker image thobe/sokar"
	@docker build -t thobe/sokar -f ci/Dockerfile .

docker.run: sep ## Runs the sokar docker image.
	@echo "--> Run docker image $(docker_image)"
	@docker run --rm --name=sokar -p 11000:11000 $(docker_image) --scale-object.name="fail-service" --saa.scale-alerts="AlertA:1.0:An upscaling alert;AlertB:-1.5:A downscaling alert"

docker.push: sep ## Pushes the sokar docker image to docker-hub
	@echo "--> Tag image to thobe/sokar:$(tag)"
	@docker tag thobe/sokar:latest thobe/sokar:$(tag)
	@echo "--> Push image thobe/sokar:latest"
	@docker push thobe/sokar:latest 
	@echo "--> Push image thobe/sokar:$(tag)"
	@docker push thobe/sokar:$(tag)

monitoring.up: ## Starts up a prometheus and a grafana instance, scraping metrics of sokar and providing a dashboard for sokar.
	make -C examples/monitoring up

monitoring.down: ## Stops the monitoring setup.
	make -C examples/monitoring down

test.service: ## Run the service-test (integration test) for sokar
	make -C test/service_test test.complete-setup

test.integration: ## Run the integration-test for sokar
	make -C test/integration test.complete-setup

lint: ## Runs the linter to check for coding-style issues
	@echo "--> Lint project"
	@echo "!!!!golangci-lint has to be installed. See: https://github.com/golangci/golangci-lint#install"
	@golangci-lint run --fast

report.test: sep ## Runs all unittests and generates a coverage- and a test-report.
	@echo "--> Run the unit-tests"	
	@go test ${packages} -timeout 30s -covermode=atomic -coverprofile=coverage.out -json | tee test-report.out

report.lint: ## Runs the linter to check for coding-style issues and generates the report file used in the ci pipeline
	@echo "--> Lint project + Reporting"
	@echo "!!!!golangci-lint has to be installed. See: https://github.com/golangci/golangci-lint#install"
	@golangci-lint run --fast --out-format checkstyle | tee lint.out

sep:
	@echo "----------------------------------------------------------------------------------"

finish:
	@echo "=================================================================================="
