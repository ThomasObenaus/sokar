.DEFAULT_GOAL				:= all
name 								:= "test_app"
build_destination := "."

all: run test stop

# This target (taken from: https://gist.github.com/prwhite/8168133) is an easy way to print out a usage/ help of all make targets.
# For all make targets the text after \#\# will be printed.
help: ## Prints the help
	@echo "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\1\:\2/' | column -c2 -t -s :)"

build: ## builds the app.
	@echo "--> Build the $(name) in $(build_destination)"
	@go build -v -o $(name) .

run: build ## Builds + runs the app
	@echo "--> Run $(name)"
	@./$(name) &

stop: ## Stops/ Kills the app
	@echo "--> Stop $(name)"
	@pkill $(name)

deps-install: ## Install the dependencies.
	@echo "--> install dependencies as listed in Gopkg.toml and Gopkg.lock"
	@dep ensure -v

test: ## runs the test and starts the app beforehand
	@echo "--> Testing $(name)"
	go test -v .

up:
	@echo "--> Startup (+build) monitoring components"
	@docker-compose up --build -d