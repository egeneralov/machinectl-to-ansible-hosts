package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type MachineRaw struct {
	Addresses string `json:"addresses,omitempty"`
	Class     string `json:"class,omitempty"`
	Machine   string `json:"machine,omitempty"`
	Os        string `json:"os,omitempty"`
	Service   string `json:"service,omitempty"`
	Version   string `json:"version,omitempty"`
}

type Machine struct {
	Name      string `json:"name"`
	Addresses string `json:"addresses"`
}

func List() ([]MachineRaw, error) {
	cmd := exec.Command("machinectl", "list", "-o", "json")
	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var machines []MachineRaw
	err = json.Unmarshal(stdout, &machines)
	if err != nil {
		return nil, err
	}
	return machines, nil
}

func MachinesRawToMachines(m []MachineRaw) []Machine {
	var machines []Machine
	for _, machine := range m {
		if machine.Class != "container" || machine.Service != "systemd-nspawn" || machine.Addresses == "" {
			continue
		}
		var ipv4 string
		for _, address := range strings.Split(machine.Addresses, "\n") {
			if strings.Count(address, ".") == 3 {
				ipv4 = address
				break
			}
		}
		if ipv4 != "" && machine.Machine != "" {
			machines = append(machines, Machine{
				Addresses: ipv4,
				Name:      machine.Machine,
			})
		}
	}
	return machines
}

func main() {
	if raw, err := List(); err != nil {
		panic(err)
	} else {
		machines := MachinesRawToMachines(raw)
		if j, err := json.Marshal(machines); err != nil {
			fmt.Println(string(j))
		}
	}
}
