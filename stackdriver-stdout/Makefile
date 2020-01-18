include ../vars.mk

.PHONY: build
build:
	yarn compile

.PHONY: deploy
deploy: build
	$(GCLOUD) app deploy
