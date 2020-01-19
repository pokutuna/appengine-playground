include ../vars.mk

go.mod go.sum main.go:
	cp ../error-reporting-go/$@ .

.PHONY: deploy
deploy: go.mod go.sum main.go
	$(GCLOUD) app deploy

clean:
	$(RM) go.mod go.sum main.goma
