.PHONY: help

LISTEN_ADDR = 0.0.0.0:8000



help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-17s\033[0m %s\n", $$1, $$2}'



test:  ## run tests
	curl localhost:8080/featured-problem

stub: ## build iprover stub
	go build -o iprover-stub/iprover iprover-stub/main.go


