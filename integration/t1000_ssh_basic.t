#!/bin/sh

test_description="Basic ssh functionality"

. ./sharness.sh

test_expect_success "We can make an SSH connection" "
	herd -l DEBUG run ssh.example.com -- uptime | grep load
"

test_done
