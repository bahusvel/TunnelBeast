package iptables

import (
	"errors"
	"log"
	"os/exec"
)

type NATEntry struct {
	SourceIP      string
	DestinationIP string
	ExternalPort  string
	InternalPort  string
}

var INTERFACE = "eth0"

func Init(Interface string) error {
	INTERFACE = Interface
	err := exec.Command("sysctl", "-w", "net.ipv4.ip_forward=1").Run()
	if err != nil {
		return errors.New("Sysctl " + err.Error())
	}
	err = exec.Command("iptables", "-t", "nat", "--flush").Run()
	if err != nil {
		return errors.New("Flush " + err.Error())
	}
	err = exec.Command("iptables", "-t", "nat", "-A", "POSTROUTING", "-o", INTERFACE, "-j", "MASQUERADE").Run()
	if err != nil {
		return errors.New("Masquerade " + err.Error())
	}
	err = exec.Command("iptables", "-t", "nat", "-A", "PREROUTING", "-i", INTERFACE, "-p", "tcp", "--destination-port", "666", "-j", "ACCEPT").Run()
	if err != nil {
		return errors.New("Port 666 " + err.Error())
	}
	return nil
}

func NewRoute(entry NATEntry) error {
	cmd := exec.Command("iptables", "-t", "nat", "-A", "PREROUTING", "-i", INTERFACE, "-p", "tcp", "-s", entry.SourceIP, "-j", "DNAT", "--to-destination", entry.DestinationIP, "--dport", entry.ExternalPort, "--to-port", entry.InternalPort)
	data, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(data))
		return err
	}
	cmd = exec.Command("iptables", "-t", "nat", "-A", "PREROUTING", "-i", INTERFACE, "-p", "udp", "-s", entry.SourceIP, "-j", "DNAT", "--to-destination", entry.DestinationIP, "--dport", entry.ExternalPort, "--to-port", entry.InternalPort)
	data, err = cmd.CombinedOutput()
	if err != nil {
		log.Println(string(data))
		return err
	}
	return nil
}

func DeleteRoute(entry NATEntry) error {
	cmd := exec.Command("iptables", "-t", "nat", "-D", "PREROUTING", "-i", INTERFACE, "-p", "tcp", "-s", entry.SourceIP, "-j", "DNAT", "--to-destination", entry.DestinationIP, "--dport", entry.ExternalPort, "--to-port", entry.InternalPort)
	data, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(data))
		return err
	}
	cmd = exec.Command("iptables", "-t", "nat", "-D", "PREROUTING", "-i", INTERFACE, "-p", "udp", "-s", entry.SourceIP, "-j", "DNAT", "--to-destination", entry.DestinationIP, "--dport", entry.ExternalPort, "--to-port", entry.InternalPort)
	data, err = cmd.CombinedOutput()
	if err != nil {
		log.Println(string(data))
		return err
	}
	return nil
}
