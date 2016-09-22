go install
BINFILE=$(which TunnelBeast)
ssh root@192.168.1.85 "killall TunnelBeast"
scp $BINFILE root@192.168.1.85:/usr/local/bin/
scp config.yml root@192.168.1.85:./
ssh root@192.168.1.85 "TunnelBeast config.yml"
