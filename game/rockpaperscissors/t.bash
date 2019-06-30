_t_pl_req() {
    [[ ! $(which curl jq) || -z $1 ]] && return 1 

    local \
        BASE=${2:-http://localhost:3000} \
        COOKIE_NAME=plateau \
        COOKIE_FILE=$1.out \
        USERINFO="{\"username\":\"$1\",\"password\":\"$1\"}"

    curl $BASE/user/register -d $USERINFO &>/dev/null
    curl $BASE/user/login --cookie-jar $COOKIE_FILE -d $USERINFO 2>/dev/null

    local match_id=$(curl -b $COOKIE_FILE $BASE/api/matchs 2>/dev/null | jq -r '.[0]')

    if [[ $match_id == "null" ]]
    then
        match_id=$(curl 2>/dev/null -b $COOKIE_FILE -X POST $BASE/api/matchs \
            -d '{"number_of_players_required":2}' | jq -r .id)
    fi

    curl $BASE/api/matchs/$match_id$3 -b $COOKIE_FILE ${@:3} 2>/dev/null | jq .
}

t_pl_cleanup() {
    [[ -z $1 ]] && return 1

    rm -f $1.out
}

t_pl_match() {
    [[ -z $1 ]] && return 1

    _t_pl_req $1 "$2" ""
}

t_pl_deals() {
    [[ -z $1 ]] && return 1

    _t_pl_req $1 "$2" /deals
}

t_pl_send() {
    [[ -z $1 ]] && return 1

    _t_pl_req $1 "$3" / -X PATCH -d "{\"request\":\"$2\"}"
}

t_pl_setupmatch() {
    [[ -z $1 || -z $2 ]] && return 1 

    (
        t_pl_send $1 PLAYER_WANT_TO_JOIN $3 && \
        t_pl_send $2 PLAYER_WANT_TO_JOIN $3 && \
        t_pl_send $2 PLAYER_WANT_TO_START_THE_GAME $3 && \
        t_pl_send $2 PLAYER_ACCEPTS $3 && \
        t_pl_send $1 PLAYER_ACCEPTS $3
    ) 1>/dev/null && echo "Done"
}
