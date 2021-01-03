#!/bin/bash

test_description="Test the keyscan command"

. ./sharness.sh

for keytype in ssh-rsa ecdsa-sha2-nistp256 ssh-ed25519; do
    test_expect_success "We can scan $keytype keys and get only that type" "
        katyusha keyscan --type $keytype ssh.example.com > output &&
        grep $keytype output &&
        ! grep -v $keytype output"
done

test_expect_success "We can scan for keys with abrreviated types" '
    katyusha keyscan --type ecdsa ssh.example.com | grep ecdsa-sha2-nistp256 &&
    katyusha keyscan --type rsa ssh.example.com | grep ssh-rsa &&
    katyusha keyscan --type ed25519 ssh.example.com | grep ssh-ed25519
'

test_done
