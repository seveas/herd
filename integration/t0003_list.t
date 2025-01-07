#!/bin/sh

test_description="Test the list command"

. ./sharness.sh

test_expect_success "List all hosts when not specifying anything" "
    herd -l debug list >out &&
    grep t0003.example.com out
"

test_expect_success "When using file:, list all hosts in the file" "
	echo t0003-bis.example.com >file &&
	herd -l debug list file:file >out &&
	grep t0003-bis.example.com out
"

test_done
