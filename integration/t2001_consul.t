#!/bin/sh

test_description="HTTP provider"

. ./sharness.sh

test_expect_success "We see hosts using the consul provider" "
    export CONSUL_HTTP_ADDR=http://consul-server-dc1.example.com:8500 &&
    echo 'datacenter   count' > expected &&
    echo 'dc1          6    ' >> expected &&
    echo 'dc2          6    ' >> expected &&
    katyusha list katyusha_provider=consul --stats datacenter --sort datacenter > actual &&
    test_cmp expected actual
"

test_expect_success "We see services using the consul provider" "
    export CONSUL_HTTP_ADDR=http://consul-server-dc1.example.com:8500 &&
    echo 'service    count' > expected &&
    echo '[consul]   2    ' >> expected &&
    echo '<nil>      10   ' >> expected &&
    katyusha list katyusha_provider=consul --stats service --sort service > actual &&
    test_cmp expected actual
"

test_done
