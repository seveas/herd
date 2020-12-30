#!/bin/sh

test_description="Basic 'can we run' test"

. ./sharness.sh

test_expect_success "We can run katyusha" "
	katyusha --help
"

test_done
