#!/bin/sh

test_description="Basic ssh functionality"

. ./sharness.sh

test_expect_success "We can scan SSH keys" "
    mkdir ~/.ssh &&
    echo 'StrictHostKeyChecking yes' > .ssh/config &&
    ssh-keyscan ssh.example.com ssh-rsa.example.com ssh-ecdsa.example.com ssh-ed25519.example.com > ~/.ssh/known_hosts_prime
"
for keytype in rsa ecdsa ed25519; do
    test_expect_success "We can connect to an $keytype host" "
        cp ~/.ssh/known_hosts_prime ~/.ssh/known_hosts &&
        herd run ssh-$keytype.example.com -- 'uptime | grep load'
    "
done

for keytype in rsa ecdsa; do
    test_expect_failure "We can not connect to a host with all keys when we have just an $keytype key" "
        grep $keytype ~/.ssh/known_hosts_prime > ~/.ssh/known_hosts &&
        herd run ssh.example.com -- 'uptime | grep load'
    "
done

keytype=ed25519
test_expect_success "We can connect to a host with all keys when we have just an $keytype key" "
    grep $keytype ~/.ssh/known_hosts_prime > ~/.ssh/known_hosts &&
    herd run ssh.example.com -- 'uptime | grep load'
"

: > .ssh/config
for keytype in rsa ecdsa ed25519; do
    test_expect_success "We can connect to a host with all keys when we have just an $keytype key" "
        grep $keytype ~/.ssh/known_hosts_prime > ~/.ssh/known_hosts &&
        herd run ssh.example.com -- 'uptime | grep load'
    "
done

test_done
