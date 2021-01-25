#!/bin/sh

test_description="Test the list command"

. ./sharness.sh

test_expect_success "List all hosts when not specifying anything" "
    katyusha -l debug list >out &&
    grep t0003.example.com out
"

test_done
