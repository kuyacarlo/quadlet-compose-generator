NAME = quadlet-compose-generator
VERSION = $(shell grep '^Version:' ../$(NAME).spec | awk '{print $$2}')

srpm:
	cd .. && \
	mkdir -p $(NAME)-$(VERSION) && \
	cp -a main.go generate.go parse.go main_test.go go.mod go.sum vendor testdata Makefile LICENSE $(NAME).spec .gitignore $(NAME)-$(VERSION)/ && \
	tar czf $(NAME)-$(VERSION).tar.gz $(NAME)-$(VERSION) && \
	rm -rf $(NAME)-$(VERSION)
	rpmbuild -bs ../$(NAME).spec \
		--define "_sourcedir $(CURDIR)/.." \
		--define "_srcrpmdir $(outdir)"
