SHARNESS_TEST_SRCDIR=sharness
. "$SHARNESS_TEST_SRCDIR/sharness.sh"
DATADIR="$SHARNESS_TEST_DIRECTORY/$this_test"
if [ -z "$SSH_AUTH_SOCK" ]; then
    eval $(ssh-agent)
    ssh-add "$SHARNESS_TEST_DIRECTORY"/openssh/user.key
fi
