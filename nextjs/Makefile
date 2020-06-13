include ../vars.mk

.PHONY: build
build:
	yarn build

.PHONY: deploy
deploy: build
	$(GCLOUD) app deploy
