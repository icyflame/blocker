#!/usr/bin/env sh

wget -O list.oisd.nl --header 'accept: text/plain' https://big.oisd.nl/hosts
wget -O list.stevenblack.hosts --header 'accept: test/plain' https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts


echo "Hosts from OISD.NL: $(wc -l list.oisd.nl)"
rm -f list.oisd.nl.cleaned
cat list.oisd.nl | rg '^0.0.0.0' > list.oisd.nl.cleaned
echo >> list.oisd.nl.cleaned
echo "Hosts from OISD.NL (after cleanup): $(wc -l list.oisd.nl.cleaned)"
rm list.oisd.nl

echo "Hosts from stevenblack: $(wc -l list.stevenblack.hosts)"
rm -f list.stevenblack.hosts.cleaned
cat list.stevenblack.hosts | rg '^0.0.0.0' > list.stevenblack.hosts.cleaned
echo >> list.stevenblack.hosts.cleaned
echo "Hosts from stevenblack (after cleanup): $(wc -l list.stevenblack.hosts.cleaned)"
rm list.stevenblack.hosts

rm -f list.blocklist-in-full
sort -u list.oisd.nl.cleaned list.stevenblack.hosts.cleaned > list.blocklist-in-full

rm list.oisd.nl.cleaned list.stevenblack.hosts.cleaned
