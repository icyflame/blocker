#!/usr/bin/env sh

wget -O list.oisd.nl --header 'accept: text/plain' https://big.oisd.nl
wget -O list.stevenblack.hosts --header 'accept: test/plain' https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts
rm -f list.blocklist-in-full
cat list.oisd.nl list.stevenblack.hosts | rg -v '^(\s+)?#' | sort -u | rg '^0.0.0.0' > list.blocklist-in-full
rm list.oisd.nl list.stevenblack.hosts
