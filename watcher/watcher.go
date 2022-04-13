package watcher

import (
	"log"

	"github.com/gosnmp/gosnmp"
)

type Watcher struct {
	Locations map[string]Location
}

type Location struct {
	Address string
	Conn    *gosnmp.GoSNMP
	Sockets []Socket
}

type Socket struct {
	Name    string
	Warning string
	Address string
}

func NewWatcher(l map[string]Location) (*Watcher, error) {
	return &Watcher{
		Locations: l,
	}, nil
}

func (w *Watcher) Walk() error {
	for _, v := range w.Locations {
		var oids []string
		var names []string
		for _, vv := range v.Sockets {
			oids = append(oids, vv.Address)
			names = append(names, vv.Name)
		}
		result, err := v.Conn.Get(oids)
		if err != nil {
			return err
		}
		for i, vv := range result.Variables {
			log.Printf("The socket with name: %s, has value of: %d.\n", names[i], vv.Value)
		}
	}
	return nil
}

func (w *Watcher) ToggleSocket(locationName, socketName string, value int) error {
	l := w.Locations[locationName]
	s := Socket{}
	for _, v := range l.Sockets {
		if v.Name == socketName {
			s = v
			break
		}
	}
	setArgs := []gosnmp.SnmpPDU{{Name: s.Address, Value: value, Type: gosnmp.Integer}}

	result, err := l.Conn.Set(setArgs)
	if err != nil {
		return err
	}
	for _, vv := range result.Variables {
		log.Printf("The socket with name: %s, has value of: %d.\n", socketName, vv.Value)
	}
	return nil
}
