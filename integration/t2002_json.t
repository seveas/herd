#!/bin/sh

test_description="JSON provider"

. ./sharness.sh

test_expect_success "We see hosts using the json provider" "
    echo test-1.example.com > expected &&
    herd list herd_provider=inventory > actual &&
    test_cmp expected actual
"

test_done
