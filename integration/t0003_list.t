#!/bin/sh

test_description="Test the list command"

. ./sharness.sh

test_expect_success "List all hosts when not specifying anything" "
    herd list >out &&
    grep t0003.example.com out
"

test_expect_success "When using file:, list all hosts in the file" "
	echo t0003-bis.example.com >file &&
	herd -l debug list file:file >out &&
	grep t0003-bis.example.com out
"

test_expect_success "It supports wildcard attributes" "
    herd list --attributes key* >out 2>&1 &&
    grep t0003.example.com out &&
    grep key:1 out &&
    grep key:2 out
"

test_done
