.DEFAULT_GOAL				:= all
name 								:= "sokar-bin"

all: build test tools cover finish

help:
	@echo "Available make targets:"
	@echo "\t- run\t\t\tBuilds + runs sokar locally."
	@echo "\t- build\t\t\tBuilds the sokar binary."
	@echo "\t- monitoring.start\tStarts up a prometheus and a grafana instance,"
	@echo "\t\t\t\tscraping metrics of sokar and providing a dashboard for sokar."
	@echo "\t- test\t\t\tRuns all unittests."
	@echo "\t- cover\t\t\tRuns the unittests and generates a coverage report."
	@echo "\t- cover.upload\t\tUploads the unittest coverage to coveralls"
	@echo "\t\t\t\t(for this the SOKAR_COVERALLS_REPO_TOKEN has to be set correctly)."
	@echo "\t- depend.install\tInstall the dependencies."
	@echo "\t- depend.update\t\tUpdate the installed dependencies."
	@echo "\t- tools\t\t\tInstalls needed tools (i.e. mock generators)."
	@echo "\t- generate.mocks\tGenerates test doubles (mocks)."


.PHONY: test
test: generate.mocks
	@echo "----------------------------------------------------------------------------------"
	@echo "--> Run the unit-tests"
	@go test ./config ./alertmanager ./nomad ./logging ./scaler ./helper ./scaleAlertAggregator ./sokar ./capacityPlanner -v

.PHONY: cover
cover: 
	@echo "----------------------------------------------------------------------------------"
	@echo "--> Run the unit-tests + coverage"
	@go test ./config ./alertmanager ./nomad ./logging ./scaler ./helper ./scaleAlertAggregator ./sokar ./capacityPlanner -v -covermode=count -coverprofile=coverage.out

cover.upload:
	# for this to get working you have to export the repo_token for your repo at coveralls.io
	# i.e. export SOKAR_COVERALLS_REPO_TOKEN=<your token>
	@${GOPATH}/bin/goveralls -coverprofile=coverage.out -service=circleci -repotoken=${SOKAR_COVERALLS_REPO_TOKEN}

#-----------------
#-- build
#-----------------
build:
	@echo "----------------------------------------------------------------------------------"
	@echo "--> Build the $(name)"
	@go build -o $(name) .

#------------------
#-- dependencies
#------------------
depend.update:
	@echo "----------------------------------------------------------------------------------"
	@echo "--> updating dependencies from Gopkg.lock"
	@dep ensure -update -v

depend.install:
	@echo "----------------------------------------------------------------------------------"
	@echo "--> install dependencies as listed in Gopkg.toml"
	@dep ensure -v

#------------------
#-- Tools
#------------------
tools:
	@go get golang.org/x/tools/cmd/cover
	@go get github.com/mattn/goveralls

generate.mocks:
	@echo "----------------------------------------------------------------------------------"
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

vendor: depend.install depend.update

run: build
	@echo "----------------------------------------------------------------------------------"
	@echo "--> Run ${name}"
	./${name} --config-file="examples/config/full.yaml" --nomad-server-address="http://${LOCAL_IP}:4646"
	# --oneshot

monitoring.start:
	@echo "----------------------------------------------------------------------------------"
	@echo "--> Startup (+build) monitoring components"
	@cd examples/monitoring && docker-compose up --build -d
	@xdg-open http://localhost:3000

monitoring.stop:
	@echo "----------------------------------------------------------------------------------"
	@echo "--> Stop monitoring components"
	@cd examples/monitoring && docker-compose down

finish:
	@echo "=================================================================================="