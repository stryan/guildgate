
#set variables
GIT_COMMIT := $(shell git rev-list -1 HEAD)
ifeq ($(PREFIX),) # PREFIX is environment variable, but if it is not set, then set default value
	PREFIX := /usr/local
endif

guildgate:
	go build -ldflags "-X main.GitCommit=$(GIT_COMMIT)"

clean: 
	rm -f guildgate

install:
	install -d $(DESTDIR)$(PREFIX)/bin
	install -m 755 guildgate $(DESTDIR)$(PREFIX)/bin
	install -d $(DESTDIR)$(PREFIX)/share/guildgate/templates
	install -D templates/* $(DESTDIR)$(PREFIX)/share/guildgate/templates

