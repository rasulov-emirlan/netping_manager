package rest

import "github.com/rasulov-emirlan/netping-manager/internal/manager"

func toServiceLocation(
	name, realLocation, SNMPaddress, SNMPcommunity string,
	SNMPport int) manager.Location {
	return manager.Location{
		Name:          name,
		RealLocation:  realLocation,
		SNMPaddress:   SNMPaddress,
		SNMPport:      SNMPport,
		SNMPcommunity: SNMPcommunity,
	}
}

func toServiceSocket(name, SNMPmib string, objectType int) manager.Socket {
	return manager.Socket{
		Name:       name,
		SNMPmib:    SNMPmib,
		ObjectType: objectType,
	}
}
