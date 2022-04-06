#/bin/bash

ipset -! create blocked-ips hash:ip
iptables -I INPUT -p tcp -m multiport --dports 80,443 -m set --match-set blocked-ips src -j DROP
./udp 127.0.0.1:8080
