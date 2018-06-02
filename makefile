all: get_deps install doc

GO=$(shell which go)
GOGET=$(GO) get

clean:
	-rm xpipe

get_deps:
	@echo -n "get dependencies... "
	@$(GOGET) github.com/moovweb/gokogiri
	@echo ok

install:
	go install ./xpipe

doc:
	bin/makedoc.rb > PROCESSES.md

.PHONY: install
