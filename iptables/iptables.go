package iptables

import (
	"log"
	"os/exec"
)

const INTERFACE = "eth0"

func Init() error {
	err := exec.Command("sysctl", "-w", "net.ipv4.ip_forward=1").Run()
	if err != nil {
		return err
	}
	err = exec.Command("iptables", "-t", "nat", "--flush").Run()
	if err != nil {
		return err
	}
	err = exec.Command("iptables", "-t", "nat", "-A", "POSTROUTING", "-o", INTERFACE, "-j", "MASQUERADE").Run()
	if err != nil {
		return err
	}
	err = exec.Command("iptables", "-t", "nat", "-A", "PREROUTING", "-i", INTERFACE, "-p", "tcp", "--destination-port", "666", "-j", "ACCEPT").Run()
	if err != nil {
		return err
	}
	return nil
}

func NewRoute(srcip string, dstip string) error {
	cmd := exec.Command("iptables", "-t", "nat", "-A", "PREROUTING", "-i", INTERFACE, "-s", srcip, "-j", "DNAT", "--to-destination", dstip)
	data, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(data))
		return err
	}
	return nil
}

func DeleteRoute(srcip string, dstip string) error {
	cmd := exec.Command("iptables", "-t", "nat", "-D", "PREROUTING", "-i", INTERFACE, "-s", srcip, "-j", "DNAT", "--to-destination", dstip)
	data, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(data))
		return err
	}
	return nil
}
