ifeq ($(shell uname -s), Darwin)
# Need to special-case ncurses pkgconfig
PKG_CONFIG_PATH:=/usr/local/opt/ncurses/lib/pkgconfig:${PKG_CONFIG_PATH}
endif

fmt:
	for d in . internal/*; do go fmt $$d/*.go; done

install:
	go install litebrew.go
