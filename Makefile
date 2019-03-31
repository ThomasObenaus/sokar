.DEFAULT_GOAL				:= all
name 								:= "sokar-bin"

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
	@go test ./config ./alertmanager ./nomad ./logging ./scaler ./helper ./scaleAlertAggregator ./sokar ./capacityPlanner ./ -v -covermode=count -coverprofile=coverage.out

cover-upload: sep
	# for this to get working you have to export the repo_token for your repo at coveralls.io
	# i.e. export SOKAR_COVERALLS_REPO_TOKEN=<your token>
	@${GOPATH}/bin/goveralls -coverprofile=coverage.out -service=circleci -repotoken=${SOKAR_COVERALLS_REPO_TOKEN}

build: sep
	@echo "--> Build the $(name)"
	@go build -o $(name) .

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
	@mockgen -source=scaler/scalingtarget_IF.go -destination test/scaler/mock_scalingtarget_IF.go 
	@mockgen -source=sokar/iface/scaler_IF.go -destination test/sokar/mock_scaler_IF.go 
	@mockgen -source=sokar/iface/capacity_planner_IF.go -destination test/sokar/mock_capacity_planner_IF.go 
	@mockgen -source=sokar/iface/scaleEventEmitter_IF.go -destination test/sokar/mock_scaleEventEmitter_IF.go 
	@mockgen -source=metrics/metrics.go -destination test/metrics/mock_metrics.go 
	@mockgen -source=logging/loggerfactory.go -destination test/logging/mock_logging.go
	@mockgen -source=runnable.go -destination test/mock_runnable.go

run: sep build
	@echo "--> Run ${name}"
	./${name} --config-file="examples/config/minimal.yaml" --nomad-server-address="http://${LOCAL_IP}:4646" --dry-run

monitoring.up:
	make -C examples/monitoring up

monitoring.down:
	make -C examples/monitoring down

sep:
	@echo "----------------------------------------------------------------------------------"

finish:
	@echo "=================================================================================="
