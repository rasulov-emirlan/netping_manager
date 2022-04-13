package watcher

import (
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

type walkResponse struct {
	Values map[string]int `json:"values"`
}

func (w *Watcher) Walk() (*walkResponse, error) {
	resp := &walkResponse{
		Values: make(map[string]int),
	}
	for _, v := range w.Locations {
		var oids []string
		var names []string
		for _, vv := range v.Sockets {
			oids = append(oids, vv.Address)
			names = append(names, vv.Name)
		}
		result, err := v.Conn.Get(oids)
		if err != nil {
			return nil, err
		}
		for i, vv := range result.Variables {
			resp.Values[names[i]] = vv.Value.(int)
		}
	}
	return resp, nil
}

type toggleSocketResponse struct {
	Values map[string]int `json:"values"`
}

func (w *Watcher) ToggleSocket(locationName, socketName string, value int) (*toggleSocketResponse, error) {
	resp := &toggleSocketResponse{
		Values: make(map[string]int),
	}
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
		return nil, err
	}
	for _, vv := range result.Variables {
		resp.Values[socketName] = vv.Value.(int)
	}
	return resp, nil
}
