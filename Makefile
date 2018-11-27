.PHONY: prod_deploy_all deploy_all

project = lambdalogger
func_name = $(project)_dev
bin = $(project).out
zip = $(project)-latest.zip

REGION ?= us-east-1
BUCKET ?= netlify-infrastructure
PROFILE ?= default
ARCH ?= amd64
OS ?= linux

sha = $(shell git rev-parse HEAD | cut -b -6)
tag = $(shell git show-ref --tags -d | grep $(sha) | cut -d '/' -f 3-)
ldflags = -X github.com/netlify/$(project)/version.SHA=$(sha) -X github.com/netlify/$(project)/version.Tag=$(tag)
clean = $(shell git status --porcelain)

###################################################################################################
# CLEAN
###################################################################################################
clean:
	rm -f dist/*.out dist/*.zip

###################################################################################################
# BUILD
###################################################################################################
build: *.go
	@echo "=== Building '$(bin)' sha: '$(sha)' tag: '$(tag)'"
	GOOS=$(OS) GOARCH=$(arch) go build -o dist/$(bin) -ldflags "$(ldflags)"

###################################################################################################
# DEPLOY
###################################################################################################
deploy: clean upload
	@echo "=== Updating function: $(func_name)"
	@aws --profile $(PROFILE) lambda update-function-code --s3-bucket $(BUCKET) --s3-key lambda/$(project)/$(zip) --function-name $(func_name) --region $(REGION) > /dev/null
	@echo "=== Finished pushing to DEV env"

deploy_test:
	cd test_func && $(MAKE) deploy
upload_test:
	cd test_func && $(MAKE) upload

upload: make_zip
	@echo "=== Uploading $(BUCKET)/lambda/$(project)/$(zip)"
	@aws --profile $(PROFILE) s3 cp dist/$(zip) s3://$(BUCKET)/lambda/$(project)/$(zip)

make_zip: build
	@which aws > /dev/null || { echo "Missing aws cli"; exit 1; }
	@cd dist/ && zip -m $(zip) $(bin)

prod_deploy: zip = $(project)-$(sha).zip
prod_deploy: _check_clean clean upload
	@echo "=== Uploading $(BUCKET)/lambda/$(project)/$(zip)"

_check_clean:
	@echo "=== Checking if it is a clean repo"
	@[ "xx$(clean)xx" == "xxxx" ] || { echo "Can't deploy from a dirty git checkout"; exit 1; }

info:
	@echo profile: $(PROFILE)
	@echo region: $(REGION)
	@echo bucket: $(BUCKET)
	@echo os: $(OS)
	@echo arch: $(ARCH)
