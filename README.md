![TunnelBeast Logo](tunnelbeast.png)

# TunnelBeast
Authenticated layer 3 reverse proxy. Do you have limited number of public IP addresses (maybe even only one), but you want to run many services of it? Or maybe you are a cloud provider who wants to give access to its clients, but doesnt have a large enough address block for all of them? Then TunnelBeast is for you!

# What it does
TunnelBeast exposes an authentication portal (and HTTP API) to the world, through which the user must authenticate and select internal address to proxy into. Once done the user can access TunnelBeast IP address as if it was the target machine, and TunnelBeast will do the magic(read the source code...). You may have as many clients as you want, accessing any internal IP addresses simultaneosly without interfering with each other.

# Installing
## Requirements
- Linux machine (VM will work just fine)
- Public IP(s) to expose TunnelBeast on
- Ports: 80, 666 + any ports used for services to be tunneled opened for the PublicIP (Check your firewall/router)
- Some techical knowledge (I will not explain everything)

## Guide
1. Grab a binary from [here](https://github.com/bahusvel/TunnelBeast/releases)
2. Put the binary in /usr/local/bin/
3. Setup the config file
4. Run TunnelBeast

## Run from Command Line
TunnelBeast supports both WebUI and command line
1. List current mapping:
curl --data "username=$YourUserName&password=$YourPassword" https://$TunnelBeastIP/list
2. List available external ports:
curl --data "username=$YourUserName&password=$YourPassword" https://$TunnelBeastIP/ports
3. Add a new mapping:
curl --data "username=$YourUserName&password=$YourPassword&internalip=$YourInternalIP&externalport=$YourExternalPort&internalport=$YourInternalPort" https://$TunnelBeastIP/add
4. Delete one existing mapping:
curl --data "username=$YourUserName&password=$YourPassword&internalip=$YourInternalIP&externalport=$YourExternalPort&internalport=$YourInternalPort&sourceip=$YourSourceIP" https://$TunnelBeastIP/delete

# Internals
TunnelBeast uses source IP to distinguish different clients hence if client A and client B have the same source IP (they are behind NAT/Gateway) they will see the same thing when trying to access TunnelBeast IP address. If you are worried about this you need to use TunnelBeast in multi IP mode. That way on successive logins clients will be given different public IPs to use.
