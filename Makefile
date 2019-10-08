.DEFAULT_GOAL				:= all
name 								:= "sokar-bin"
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
nomad_server := "http://${LOCAL_IP}:4646"

all: tools test build finish

# This target (taken from: https://gist.github.com/prwhite/8168133) is an easy way to print out a usage/ help of all make targets.
# For all make targets the text after \#\# will be printed.
help: ## Prints the help
	@echo "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\1\:\2/' | column -c2 -t -s :)"

.PHONY: test
test: sep gen-mocks ## Runs all unittests and generates a coverage report.
	@echo "--> Run the unit-tests"
	@go test ./config ./alertmanager ./nomad ./logging ./scaler ./helper ./scaleAlertAggregator ./sokar ./capacityPlanner ./aws ./awsEc2 ./nomadWorker ./api ./ -covermode=count -coverprofile=coverage.out

cover-upload: sep ## Uploads the unittest coverage to coveralls (for this the SOKAR_COVERALLS_REPO_TOKEN has to be set correctly).
	# for this to get working you have to export the repo_token for your repo at coveralls.io
	# i.e. export SOKAR_COVERALLS_REPO_TOKEN=<your token>
	@${GOPATH}/bin/goveralls -coverprofile=coverage.out -service=circleci -repotoken=${SOKAR_COVERALLS_REPO_TOKEN}

build: sep ## Builds the sokar binary.
	@echo "--> Build the $(name) in $(build_destination)"
	@go build -v -ldflags "-X main.version=$(tag) -X main.buildTime=$(build_time) -X main.revision=$(revision) -X main.branch=$(branch)" -o $(sokar_file_name) .

deps-update: sep ## Update the installed dependencies.
	@echo "--> updating dependencies. Trying to find newer versions as they are listed in Gopkg.lock"
	@dep ensure -update -v

deps-install: sep ## Install the dependencies, without looking for new versions of dependencies.
	@echo "--> install dependencies as listed in Gopkg.toml and Gopkg.lock"
	@dep ensure -vendor-only -v

tools: sep ## Installs needed tools (i.e. mock generators).
	@echo "--> Install needed tools."
	@go get golang.org/x/tools/cmd/cover
	@go get github.com/mattn/goveralls

gen-mocks: sep ## Generates test doubles (mocks).
	@echo "--> generate mocks (github.com/golang/mock/gomock is required for this)"
	@go get github.com/golang/mock/gomock
	@go install github.com/golang/mock/mockgen
	@mockgen -source=nomad/nomadclient_IF.go -destination test/nomad/mock_nomadclient_IF.go
	@mockgen -source=nomadWorker/nomadclient_IF.go -destination test/nomadWorker/mock_nomadclient_IF.go
	@mockgen -source=aws/iface/autoscaling_IF.go -destination test/aws/mock_autoscaling_IF.go 
	@mockgen -source=scaler/scalingtarget_IF.go -destination test/scaler/mock_scalingtarget_IF.go 
	@mockgen -source=sokar/iface/scaler_IF.go -destination test/sokar/mock_scaler_IF.go 
	@mockgen -source=sokar/iface/capacity_planner_IF.go -destination test/sokar/mock_capacity_planner_IF.go 
	@mockgen -source=sokar/iface/scaleEventEmitter_IF.go -destination test/sokar/mock_scaleEventEmitter_IF.go 
	@mockgen -source=metrics/metrics.go -destination test/metrics/mock_metrics.go 
	@mockgen -source=logging/loggerfactory.go -destination test/logging/mock_logging.go
	@mockgen -source=runnable.go -destination test/mock_runnable.go

gen-metrics-md: sep ## Generate metrics documentation (Metrics.md) based on defined metrics in code.
	@echo "--> generate Metrics.md"
	@python3 gen_metrics_doc.py > Metrics.md

run.aws-ec2: sep build ## Builds + runs sokar locally in aws ec2 mode.
	@echo "--> Run $(sokar_file_name)"
	$(sokar_file_name) --config-file="examples/config/full.yaml" --sca.mode="aws-ec2"

run.nomad-dc: sep build ## Builds + runs sokar locally in data-center mode.
	@echo "--> Run $(sokar_file_name)"
	$(sokar_file_name) --config-file="examples/config/full.yaml" --sca.nomad.server-address=$(nomad_server) --sca.mode="nomad-dc"

run.nomad-job: sep build ## Builds + runs sokar locally in job mode.
	@echo "--> Run $(sokar_file_name)"
	$(sokar_file_name) --config-file="examples/config/full.yaml" --sca.nomad.server-address=$(nomad_server)

docker.build: sep ## Builds the sokar docker image.
	@echo "--> Build docker image thobe/sokar"
	@docker build -t thobe/sokar -f ci/Dockerfile .

docker.run: sep ## Runs the sokar docker image.
	@echo "--> Run docker image $(docker_image)"
	@docker run --rm --name=sokar -p 11000:11000 $(docker_image) --sca.nomad.server-address=$(nomad_server) --scale-object.name="fail-service" --saa.scale-alerts="AlertA:1.0:An upscaling alert;AlertB:-1.5:A downscaling alert"

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

sep:
	@echo "----------------------------------------------------------------------------------"

finish:
	@echo "=================================================================================="
