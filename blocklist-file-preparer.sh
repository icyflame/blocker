#!/usr/bin/env bash

OUTFILE=$1
if [[ -z "$OUTFILE" ]];
then
	cat <<EOF
ERROR: Must provide 1 argument, which is the output file

Example:

    $0 ~/list.blocklist-in-full
EOF
	exit 43
fi

curl --silent --header 'accept: text/plain' https://raw.githubusercontent.com/icyflame/hosts/master/hosts_abp.txt > list.icyflame.abp
curl --silent --header 'accept: text/plain' https://big.oisd.nl > list.oisd.nl.abp &
# This list is not offered in the ABP format. It blocks some domains which are already blocked by
# the Big OISD list. However, I am not sure how many of the domains in the StevenBlack list are
# blocked by OISD. I have not calculated this yet. So, for now, I will continue to merge this list
# into the OISD list.
curl --silent --header 'accept: test/plain' https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts > list.stevenblack.hosts &

wait

cat list.stevenblack.hosts | grep -E '^0.0.0.0' | sed 's#0.0.0.0 ##g' | sed 's#^#||#g' | sed 's#$#^#g' > list.stevenblack.abp

echo "Hosts from icyflame: $(wc -l list.icyflame.abp)"
echo "Hosts from oisd big: $(wc -l list.oisd.nl.abp)"
echo "Hosts from stevenblack: $(wc -l list.stevenblack.abp)"
cat list.icyflame.abp list.oisd.nl.abp list.stevenblack.abp > list.blocklist-in-full

echo "Final hosts (with duplicates): $(wc -l list.blocklist-in-full)"
echo "Final hosts (without duplicates): $(sort -u list.blocklist-in-full | wc -l)"

mv list.blocklist-in-full $OUTFILE

rm -f list.icyflame.abp list.oisd.nl.abp list.stevenblack.abp list.stevenblack.hosts
