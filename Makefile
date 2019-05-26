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

help:
	@echo "Available make targets:"
	@echo "\t- run\t\t\tBuilds + runs sokar locally."
	@echo "\t- build\t\t\tBuilds the sokar binary."
	@echo "\t- monitoring.up\tStarts up a prometheus and a grafana instance,"
	@echo "\t\t\t\tscraping metrics of sokar and providing a dashboard for sokar."
	@echo "\t- test\t\t\tRuns all unittests and generates a coverage report."
	@echo "\t- cover-upload\t\tUploads the unittest coverage to coveralls"
	@echo "\t\t\t\t(for this the SOKAR_COVERALLS_REPO_TOKEN has to be set correctly)."
	@echo "\t- deps-install\tInstall the dependencies."
	@echo "\t- deps-update\t\tUpdate the installed dependencies."
	@echo "\t- tools\t\t\tInstalls needed tools (i.e. mock generators)."
	@echo "\t- gen-mocks\tGenerates test doubles (mocks)."


.PHONY: test
test: sep gen-mocks
	@echo "--> Run the unit-tests"
	@go test ./config ./alertmanager ./nomad ./logging ./scaler ./helper ./scaleAlertAggregator ./sokar ./capacityPlanner ./nomadWorker ./ -covermode=count -coverprofile=coverage.out

cover-upload: sep
	# for this to get working you have to export the repo_token for your repo at coveralls.io
	# i.e. export SOKAR_COVERALLS_REPO_TOKEN=<your token>
	@${GOPATH}/bin/goveralls -coverprofile=coverage.out -service=circleci -repotoken=${SOKAR_COVERALLS_REPO_TOKEN}

build: sep
	@echo "--> Build the $(name) in $(build_destination)"
	@go build -v -ldflags "-X main.version=$(tag) -X main.buildTime=$(build_time) -X main.revision=$(revision) -X main.branch=$(branch)" -o $(sokar_file_name) .

deps-update: sep
	@echo "--> updating dependencies. Trying to find newer versions as they are listed in Gopkg.lock"
	@dep ensure -update -v

deps-install: sep
	@echo "--> install dependencies as listed in Gopkg.toml and Gopkg.lock"
	@dep ensure -v

tools: sep
	@echo "--> Install needed tools."
	@go get golang.org/x/tools/cmd/cover
	@go get github.com/mattn/goveralls

gen-mocks: sep
	@echo "--> generate mocks (github.com/golang/mock/gomock is required for this)"
	@go get github.com/golang/mock/gomock
	@go install github.com/golang/mock/mockgen
	@mockgen -source=nomad/nomadclient_IF.go -destination test/nomad/mock_nomadclient_IF.go 
	@mockgen -source=nomadWorker/iface/autoscaling_IF.go -destination test/nomadWorker/mock_autoscaling_IF.go 
	@mockgen -source=scaler/scalingtarget_IF.go -destination test/scaler/mock_scalingtarget_IF.go 
	@mockgen -source=sokar/iface/scaler_IF.go -destination test/sokar/mock_scaler_IF.go 
	@mockgen -source=sokar/iface/capacity_planner_IF.go -destination test/sokar/mock_capacity_planner_IF.go 
	@mockgen -source=sokar/iface/scaleEventEmitter_IF.go -destination test/sokar/mock_scaleEventEmitter_IF.go 
	@mockgen -source=metrics/metrics.go -destination test/metrics/mock_metrics.go 
	@mockgen -source=logging/loggerfactory.go -destination test/logging/mock_logging.go
	@mockgen -source=runnable.go -destination test/mock_runnable.go

gen-metrics-md: sep
	@echo "--> generate Metrics.md"
	@python3 gen_metrics_doc.py > Metrics.md


run: sep build
	@echo "--> Run $(sokar_file_name)"
	$(sokar_file_name) --config-file="examples/config/full.yaml" --nomad.server-address=$(nomad_server)

docker.build: sep
	@echo "--> Build docker image thobe/sokar"
	@docker build -t thobe/sokar -f ci/Dockerfile .

docker.run: sep
	@echo "--> Run docker image $(docker_image)"
	@docker run --rm --name=sokar -p 11000:11000 $(docker_image) --nomad.server-address=$(nomad_server)

docker.push: sep
	@echo "--> Tag image to thobe/sokar:$(tag)"
	@docker tag thobe/sokar:latest thobe/sokar:$(tag)
	@echo "--> Push image thobe/sokar:latest"
	@docker push thobe/sokar:latest 
	@echo "--> Push image thobe/sokar:$(tag)"
	@docker push thobe/sokar:$(tag)



monitoring.up:
	make -C examples/monitoring up

monitoring.down:
	make -C examples/monitoring down

sep:
	@echo "----------------------------------------------------------------------------------"

finish:
	@echo "=================================================================================="
