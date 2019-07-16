. assert.sh || . assert

. helpers.bash

assert_raises "new_plateau" 1
assert_raises "new_plateau mygame"

assert_raises "test -d ~/mygame"

# $TPG_PLATEAU_BASEURL

assert_raises "tpg_setupmatch" 1
assert_raises "tpg_setupmatch P1 P2"

assert_raises "tpg_match" 1
assert_raises "tpg_match P1"
assert_raises "tpg_match P2"

assert_raises "tpg_send" 1
assert_raises "tpg_send P1" 1
assert_raises "tpg_send P1 ?"
assert_raises "tpg_send P2 ?"

assert_raises "tpg_deals" 1
assert_raises "tpg_deals P1"
assert_raises "tpg_deals P2"

assert_raises "test -f P1.cookie"
assert_raises "test -f P2.cookie"

assert_raises "tpg_cleanup" 1
assert_raises "tpg_cleanup P1"
assert_raises "tpg_cleanup P2"

assert_raises "test -f P1.cookie" 1
assert_raises "test -f P2.cookie" 1

assert_end helpers
