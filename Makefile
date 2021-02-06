ifeq ($(shell uname -s), Darwin)
# Need to special-case ncurses pkgconfig
PKG_CONFIG_PATH:=/usr/local/opt/ncurses/lib/pkgconfig:${PKG_CONFIG_PATH}
endif

fmt:
	scripts/fmt.sh

install:
	scripts/build.sh

build_all:
	scripts/build.sh all
