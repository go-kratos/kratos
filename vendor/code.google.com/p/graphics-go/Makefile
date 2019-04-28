# Copyright 2011 The Graphics-Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.inc

all: install

# The order matters: earlier packages may not depend on later ones.
DIRS=\
	graphics

test:
	cd graphics; GOMAXPROCS=2 gomake test

bench:
	cd graphics; GOMAXPROCS=2 gomake bench

install clean nuke:
	for dir in $(DIRS); do \
		$(MAKE) -C $$dir $@ || exit 1; \
	done
