#!/bin/sh

test_description="Basic 'can we run' test"

. ./sharness.sh

test_expect_success "We can run herd" "
	herd --help
"

test_done
