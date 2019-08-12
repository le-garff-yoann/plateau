. assert.sh || . assert

. helpers.bash

assert_raises "new_plateau" 1
assert_raises "new_plateau mygame"

assert_raises "test -d ~/mygame"

assert_end helpers
