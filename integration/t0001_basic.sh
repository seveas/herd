#!/bin/sh

test_description="Basic 'can we run' test"

. ./sharness.sh

test_expect_success "We can run herd" "
	herd --help
"

test_expect_success "All providers marked the base providers as squashable" "
	! git grep '^\s*baseProvider$'
"

test_done
