GOPATH := $(shell go env | grep GOPATH | sed 's/GOPATH="\(.*\)"/\1/')
PATH := $(GOPATH)/bin:$(PATH)
export $(PATH)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COPY_FILES := $(wildcard ./dist/terraform-provider-gitfile*/*)
TEST_DIR := test/
TEST_DESTS := $(patsubst $(TEST_DIR)%/test.sh, $(TEST_DIR)%, $(wildcard $(TEST_DIR)*/test.sh))

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

fetch: ## download makefile dependencies
	test -f $(GOPATH)/bin/goreleaser || go get -u -v github.com/goreleaser/goreleaser

clean: ## cleans previously built binaries
	rm -rf ./dist
	@cd test && $(MAKE) clean

publish: clean fetch ## publishes assets
	@if [ "${GITHUB_TOKEN}" == "" ]; then\
	  echo "GITHUB_TOKEN is not set";\
		exit 1;\
	fi
	@if [ "$(GIT_BRANCH)" != "master" ]; then\
	  echo "Current branch is: '$(GIT_BRANCH)'.  Please publish from 'master'";\
		exit 1;\
	fi
	git tag -a $(VERSION) -m "$(MESSAGE)"
	git push --follow-tags
	$(GOPATH)/bin/goreleaser

build: clean fetch ## publishes in dry run mode
	$(GOPATH)/bin/goreleaser --skip-publish --snapshot


.PHONY: test copyplugins

copyplugins: ## copy plugins to test folders
	$(eval OS_ARCH := $(patsubst ./dist/terraform-provider-gitfile_%/terraform-provider-gitfile, %, $(COPY_FILES)))
	$(eval TEST_FOLDERS := $(foreach p,$(OS_ARCH), $(patsubst %,%/terraform.d/plugins/$p,$(TEST_DESTS))))
	@sleep 1
	@mkdir -p $(TEST_FOLDERS)

	@for o in $(OS_ARCH); do \
	  for f in $(TEST_DESTS); do \
		  cp ./dist/terraform-provider-gitfile_$$o/* $$f/terraform.d/plugins/$$o; \
		done; \
	done

test: copyplugins ## test
	@cd test && $(MAKE) test

fulltest: build test ## build and test