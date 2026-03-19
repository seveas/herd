#!/bin/sh

test_description="Basic ssh functionality"

. ./sharness.sh

test_expect_success "We can override expected exit status" "
	herd -l DEBUG run --expect-exit-status 2 ssh.example.com -- exit 2 >out && \
	grep 'completed successfully with status 2' out
"

test_done
