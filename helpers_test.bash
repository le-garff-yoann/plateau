. assert.sh || . assert

. helpers.bash

assert_raises "new_plateau" 1
assert_raises "new_plateau mygame"

assert_raises "test -d ~/mygame"

# $T2PG_PLATEAU_BASEURL

assert_raises "t2pg_setupmatch" 1
assert_raises "t2pg_setupmatch P1" 1
assert_raises "t2pg_setupmatch P1 P2"

assert_raises "t2pg_match" 1
assert_raises "t2pg_match P1"
assert_raises "t2pg_match P2"

assert_raises "t2pg_send" 1
assert_raises "t2pg_send P1" 1
assert_raises "t2pg_send P1 ?"
assert_raises "t2pg_send P2 ?"

assert_raises "t2pg_deals" 1
assert_raises "t2pg_deals P1"
assert_raises "t2pg_deals P2"

assert_raises "test -f P1.cookie"
assert_raises "test -f P2.cookie"

assert_raises "t2pg_cleanup" 1
assert_raises "t2pg_cleanup P1"
assert_raises "t2pg_cleanup P2"

assert_raises "test -f P1.cookie" 1
assert_raises "test -f P2.cookie" 1

assert_end helpers
