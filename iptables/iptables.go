package iptables

import (
	"bytes"
	"errors"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type NATEntry struct {
	SourceIP      string
	DestinationIP string
	ExternalPort  string
	InternalPort  string
	Client        string
	Traffic       string
}

const (
	BYTE = 1.0 << (10 * iota)
	KILOBYTE
	MEGABYTE
	GIGABYTE
	TERABYTE
)

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
	err = exec.Command("iptables", "--flush").Run()
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
	cmd := exec.Command("iptables", "-t", "nat", "-A", "PREROUTING", "-i", INTERFACE, "-p", "tcp", "-s", entry.SourceIP, "-j", "DNAT", "--to-destination", entry.DestinationIP+":"+entry.InternalPort, "--dport", entry.ExternalPort)
	data, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(data))
		return err
	}
	cmd = exec.Command("iptables", "-t", "nat", "-A", "PREROUTING", "-i", INTERFACE, "-p", "udp", "-s", entry.SourceIP, "-j", "DNAT", "--to-destination", entry.DestinationIP+":"+entry.InternalPort, "--dport", entry.ExternalPort)
	data, err = cmd.CombinedOutput()
	if err != nil {
		log.Println(string(data))
		return err
	}
	cmd = exec.Command("iptables", "-A", "FORWARD", "-s", entry.DestinationIP)
	data, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}
	cmd = exec.Command("iptables", "-A", "FORWARD", "-d", entry.DestinationIP)
	data, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

func DeleteRoute(entry NATEntry) error {
	cmd := exec.Command("iptables", "-t", "nat", "-D", "PREROUTING", "-i", INTERFACE, "-p", "tcp", "-s", entry.SourceIP, "-j", "DNAT", "--to-destination", entry.DestinationIP+":"+entry.InternalPort, "--dport", entry.ExternalPort)
	data, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(data))
		return err
	}
	cmd = exec.Command("iptables", "-t", "nat", "-D", "PREROUTING", "-i", INTERFACE, "-p", "udp", "-s", entry.SourceIP, "-j", "DNAT", "--to-destination", entry.DestinationIP+":"+entry.InternalPort, "--dport", entry.ExternalPort)
	data, err = cmd.CombinedOutput()
	if err != nil {
		log.Println(string(data))
		return err
	}
	cmd = exec.Command("iptables", "-D", "FORWARD", "-s", entry.DestinationIP)
	data, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}
	cmd = exec.Command("iptables", "-D", "FORWARD", "-d", entry.DestinationIP)
	data, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

func Traffic(entry NATEntry) string {
	args := "iptables -nxvL | grep " + entry.DestinationIP + " | awk '{print $2}'"
	cmd := exec.Command("/bin/sh", "-c", args)
	data, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(data)[0])
		return "Err"
	}

	var bandwidth int64
	values := strings.Split(string(bytes.Trim(data, "\n")), "\n")
	for _, v := range values {
		b, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			log.Println(err)
			return "Err"
		}
		bandwidth = bandwidth + b
	}

	unit := ""
	value := int64(bandwidth)
	switch {
	case bandwidth >= TERABYTE:
		unit = "TB"
		value = value / TERABYTE
	case bandwidth >= GIGABYTE:
		unit = "GB"
		value = value / GIGABYTE
	case bandwidth >= MEGABYTE:
		unit = "MB"
		value = value / MEGABYTE
	case bandwidth >= KILOBYTE:
		unit = "KB"
		value = value / KILOBYTE
	case bandwidth >= BYTE:
		unit = "B"
	case bandwidth == 0:
		unit = "B"
		return "0B"
	}

	return strconv.FormatInt(value, 10) + unit
}
