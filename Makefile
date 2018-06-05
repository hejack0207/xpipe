all: get_deps install doc

GO=$(shell which go)
GOGET=$(GO) get
VERSION=v0.0.1)
DESCRIPTION=version 0.0.1

clean:
	rm -rf dist

get_deps:
	@echo -n "get dependencies... "
	@$(GOGET) github.com/moovweb/gokogiri
	@echo ok

install:
	go install

gorelease:
	git commit -m "$(DESCRIPTION)"
	git tag -d $(VERSION) -m $(DESCRIPTION)
	git tag -a $(VERSION) -m $(DESCRIPTION)
	goreleaser

doc:
	docs/makedoc.rb > PROCESSES.md

.PHONY: install
