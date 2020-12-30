#!/bin/sh

test_description="Output format/detail tests"

. ./sharness.sh

test_expect_success "Debug output shows up" "
	herd -l debug list '*' | grep 'hosts returned from'
"

test_done
