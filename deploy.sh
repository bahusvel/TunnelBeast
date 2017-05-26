go install
BINFILE=$(which TunnelBeast)
DESTINATION='192.168.1.91'

ssh root@$DESTINATION "killall TunnelBeast"
scp $BINFILE root@$DESTINATION:/usr/local/bin/
scp config.yml root@$DESTINATION:./
ssh root@$DESTINATION "TunnelBeast config.yml"
