package docker

import (
	"bytes"
	"errors"
	"github.com/c-robinson/iplib"
	"os"
	"sync"
)

type Subnetter struct {
	// mutex
	Lock *sync.Mutex
	// RegistryPath to the network's subnet registry, a list of used subnets
	RegistryPath string
	// Subnets is a list of subnets that ips can be allocated from
	Subnets []iplib.Net4
}

var (
	ErrNoAvailableSubnet = errors.New("no available subnet")
)

func (n Subnetter) Release(subnet iplib.Net4) error {
	used, err := n.Used()
	if err != nil {
		return err
	}

	subnets := bytes.Replace(used, []byte(subnet.String()+"\n"), []byte(""), 1)

	return os.WriteFile(n.RegistryPath, subnets, 0644)
}

func (n Subnetter) Allocate(subnet iplib.Net4) error {
	// append subnet to n.RegistryPath file
	f, err := os.OpenFile(n.RegistryPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer func() {
		if err = f.Close(); err != nil {
			panic(err)
		}
	}()

	_, err = f.Write([]byte(subnet.String() + "\n"))
	if err != nil {
		return err
	}

	return nil
}

func (n Subnetter) Used() ([]byte, error) {
	if _, err := os.Stat(n.RegistryPath); os.IsNotExist(err) {
		err = os.WriteFile(n.RegistryPath, []byte(""), 0600)
		if err != nil {
			return nil, err
		}

		return []byte{}, nil
	}
	return os.ReadFile(n.RegistryPath)
}

func (n Subnetter) NextSubnet() (*iplib.Net4, error) {
	n.Lock.Lock()
	defer n.Lock.Unlock()
	used, err := n.Used()
	if err != nil {
		return nil, err
	}

	// todo: make subnets race with each other to find the first available sub-subnet
	for _, subnet := range n.Subnets {
		subnets, err := subnet.Subnet(24)
		if err != nil {
			return nil, err
		}

		for _, s := range subnets {
			if bytes.Contains(used, []byte(s.String())) {
				continue
			}

			err = n.Allocate(s)
			if err != nil {
				return nil, err
			}

			return &s, nil
		}
	}

	return nil, ErrNoAvailableSubnet
}
