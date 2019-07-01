_t2pg_req() {
    [[ ! $(which curl jq) || -z $1 ]] && return 1 

    local \
        BASE=${2:-http://localhost:3000} \
        COOKIE_NAME=plateau \
        COOKIE_FILE=$1.cookie \
        USERINFO="{\"username\":\"$1\",\"password\":\"$1\"}"

    curl $BASE/user/register -d $USERINFO &>/dev/null
    curl $BASE/user/login --cookie-jar $COOKIE_FILE -d $USERINFO 2>/dev/null

    local match_id=$(curl -b $COOKIE_FILE $BASE/api/matchs 2>/dev/null | jq -r '.[0]')

    [[ $match_id == "null" ]] && \
    match_id=$(curl 2>/dev/null -b $COOKIE_FILE -X POST $BASE/api/matchs \
        -d '{"number_of_players_required":2}' | jq -r .id)

    curl $BASE/api/matchs/$match_id$3 -b $COOKIE_FILE ${@:3} 2>/dev/null | jq .
}

t2pg_cleanup() {
    [[ -z $1 ]] && return 1

    rm -f $1.cookie
}

t2pg_deals() {
    [[ -z $1 ]] && return 1

    _t2pg_req $1 "$2" /deals
}

t2pg_send() {
    [[ -z $1 ]] && return 1

    _t2pg_req $1 "$3" / -X PATCH -d "{\"request\":\"$2\"}"
}

t2pg_setupmatch() {
    [[ -z $1 || -z $2 ]] && return 1 

    (
        t2pg_send $1 PLAYER_WANT_TO_JOIN $3 && \
        t2pg_send $2 PLAYER_WANT_TO_JOIN $3 && \
        t2pg_send $2 PLAYER_WANT_TO_START_THE_GAME $3 && \
        t2pg_send $2 PLAYER_ACCEPTS $3 && \
        t2pg_send $1 PLAYER_ACCEPTS $3
    ) 1>/dev/null && echo "Done"
}
