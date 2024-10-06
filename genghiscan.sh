#!/bin/bash

get_subdomains() {
    # rapiddns
    curl -s "https://rapiddns.io/subdomain/$1\?full\=1\#result" 2>/dev/null >"$2/subs/rapiddns.html"
    find "$2" -type f -name rapiddns.html 2>/dev/null | html-tool tags td | sort -u | grep "$1" > "$2/subs/zRapiddns"
    # anubisdb
    curl -s "https://jldc.me/anubis/subdomains/$1" 2>/dev/null | grep -Po "((http|https):\/\/)?(([\w.-]*)\.([\w]*)\.([A-z]))\w+" 2>/dev/null | sort -u 2>/dev/null >"$2/subs/zAnubis"
    # subfinder
    subfinder -silent -d "$1" -all -recursive >"$2/subs/zSubfinder"
    # assetfinder
    assetfinder --subs-only "$1" >"$2/subs/zAssetfinder"
    # findomain
    findomain -q -t "$1" >"$2/subs/zFindomain"
    # github-subdomains
    github-subdomains -raw -d $1 >"$2/subs/zGithub"
    # chaos
    chaos-client -d "$1" -k $chaos -silent >"$2/subs/zChaos"
    
    cat "$2/subs/z*"|uniq -u > "$2/subs/all"
}

get_ips() {
    cat "$2/rapiddns.html" | grep "<td><a" 2>/dev/null | cut -d '"' -f 2 2>/dev/null | cut -d '/' -f3 2>/dev/null | sed 's/\#result//g' | grep -oE '((1?[0-9][0-9]?|2[0-4][0-9]|25[0-5])\.){3}(1?[0-9][0-9]?|2[0-4][0-9]|25[0-5])' >"$2/ips.txt"
    python ~/DEV/heASNs/heASNs.py -i "$1" >"$2/heAsn.json"
    cat "$2/heAsn.json" | jq .results[0].data.asns[] | sed 's/"//g' | while read -r line; do whois -h whois.radb.net -- "-i origin $line" | grep -Eo "([0-9.]+){4}/[0-9]+" | sort -u | tee ranges.txt | mapcidr -silent | dnsx -silent -ptr -resp-only | anew "$2/subs/Zrevdns"; done
}

get_urls() {
    echo "$1" | waybackurls >"$2/urls/zWayback"
    waymore -i "$1" -mode R -l 0 -v -oR "$2/urls/responses/" -oijs -ci none
}

folder_setup() {
    mkdir "$1/subs"
    mkdir "$1/urls"
}

DOMAIN=$1
OUTDIR="$2/$DOMAIN"
ORG="$(echo $DOMAIN | awk -F"." '{ print $(NF-1) }')"
chaos="$(echo $CHAOS_KEY)"

folder_setup "$OUTDIR"
get_subdomains "$DOMAIN" "$OUTDIR"
get_ips "$ORG" "$OUTDIR"
get_urls "$DOMAIN" "$OUTDIR"
