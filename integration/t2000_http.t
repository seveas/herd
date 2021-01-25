#!/bin/sh

test_description="HTTP provider"

. ./sharness.sh

test_expect_success "We see hosts in the http provider" "
	herd -l DEBUG list provider=http | grep http.example.net
"

test_done
