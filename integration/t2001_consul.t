#!/bin/sh

test_description="HTTP provider"

. ./sharness.sh

# Wait for consul to settle
(
    set -x
    export CONSUL_HTTP_ADDR=http://consul-server-dc1.example.com:8500
    while ! consul members | grep client; do sleep 1; done
    export CONSUL_HTTP_ADDR=http://consul-server-dc1.example.com:8500
    while ! consul members | grep client; do sleep 1; done
    set +x
)

test_expect_success "We see hosts using the consul provider" "
    export CONSUL_HTTP_ADDR=http://consul-server-dc1.example.com:8500 &&
    echo 'datacenter   count' > expected &&
    echo 'dc1          6    ' >> expected &&
    echo 'dc2          6    ' >> expected &&
    herd list --refresh herd_provider=consul --count datacenter --sort datacenter > actual &&
    test_cmp expected actual
"

test_expect_success "We see services using the consul provider" "
    export CONSUL_HTTP_ADDR=http://consul-server-dc1.example.com:8500 &&
    echo 'service    count' > expected &&
    echo '[consul]   2    ' >> expected &&
    echo '<nil>      10   ' >> expected &&
    herd list herd_provider=consul --count service --sort service > actual &&
    test_cmp expected actual
"

test_done
