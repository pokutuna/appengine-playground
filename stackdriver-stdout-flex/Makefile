include ../vars.mk

.PHONY: build
build:
	$(MAKE) -C ../stackdriver-stdout build
	cp -T ../stackdriver-stdout/package.json package.json
	cp -RT ../stackdriver-stdout/build ./build

.PHONY: deploy
deploy: build
	$(GCLOUD) app deploy
