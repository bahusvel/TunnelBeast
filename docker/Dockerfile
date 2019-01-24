FROM ubuntu:latest
RUN apt-get update && apt-get -y install sysstat gawk net-tools iputils-ping iftop iptables ca-certificates
COPY config.yml /config.yml
COPY TunnelBeast /TunnelBeast
ENTRYPOINT ["./TunnelBeast", "config.yml"] 

