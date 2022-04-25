package sentry

import (
	"context"
	"errors"
	"log"

	"github.com/gosnmp/gosnmp"
	"github.com/rasulov-emirlan/netping-manager/internal/manager"
)

type Sentry struct {
}

func connect(address, community string, port int) (*gosnmp.GoSNMP, error) {
	g := *gosnmp.Default
	g.Target = address
	g.Community = community
	g.Port = uint16(port)
	if err := g.Connect(); err != nil {
		return nil, err
	}
	return &g, nil
}

// TODO: not sure if this function should incapsulate some logic from
// buisness logic. It would speed up our buisness logic if we incapsulated it
// But it does not look too good
func (s *Sentry) CheckSocket(ctx context.Context, mib []string, address, community string, port int) ([]bool, error) {
	g, err := connect(address, community, port)
	if err != nil {
		return nil, err
	}
	res, err := g.Get(mib)
	if err != nil {
		return nil, err
	}
	checks := []bool{}

	for _, v := range res.Variables {
		result, ok := v.Value.(int)
		if !ok {
			return nil, errors.New("sentry: incorrect value type")
		}
		if result == 0 {
			checks = append(checks, true)
			continue
		}
		checks = append(checks, false)
	}
	log.Println(checks)
	return checks, nil
}

func (s *Sentry) CheckSockets(ctx context.Context, oids []string, address, community string, port int) ([]*manager.Socket, error) {
	g, err := connect(address, community, port)
	if err != nil {
		return nil, err
	}
	res, err := g.Get(oids)
	if err != nil {
		return nil, err
	}
	var (
		response []*manager.Socket
	)
	for i, v := range res.Variables {
		ss := &manager.Socket{
			IsON:          false,
			SNMPaddress:   address,
			SNMPcommunity: community,
			SNMPport:      port,
			SNMPmib:       oids[i],
		}
		vv, ok := v.Value.(int)
		if !ok {
			return nil, errors.New("sentry: inorrect type of value")
		}
		if vv == 1 {
			ss.IsON = true
		}
		response = append(response, ss)
	}
	return response, nil
}

func (s *Sentry) ToggleSocket(ctx context.Context, turnOn bool, socketMIB, address, community string, port int) (*manager.Socket, error) {
	g, err := connect(address, community, port)
	if err != nil {
		return nil, err
	}
	var turnOnOrOff int = 0
	if !turnOn {
		turnOnOrOff = 1
	}
	input := []gosnmp.SnmpPDU{{
		Value: turnOnOrOff,
		Type:  gosnmp.Integer,
		Name:  socketMIB,
	}}
	_, err = g.Set(input)
	if err != nil {
		return nil, err
	}
	return &manager.Socket{
		SNMPaddress:   address,
		SNMPport:      port,
		SNMPcommunity: community,
		SNMPmib:       socketMIB,
		IsON:          turnOn,
	}, nil
}
