package manager

type Location struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	RealLocation string `json:"realLocation"`

	SNMPaddress   string `json:"snmpAddress"`
	SNMPport      int    `json:"snmpPort"`
	SNMPcommunity string `json:"snmpCommunity"`

	Sockets []*Socket `json:"sockets"`
}

type Socket struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`

	SNMPaddress   string `json:"snmpAddress"`
	SNMPport      int    `json:"snmpPort"`
	SNMPcommunity string `json:"snmpCommunity"`
	
	SNMPmib string `json:"snmpMib"`

	IsON    bool   `json:"isON"`

	ObjectType int `json:"objectType"`
}

// This is just an example what object types might be
// It is unlikely that you will use these constants
// since in future you may add more types in your database
const (
	TypeUnknown = iota + 1
	TypeAC
	TypeHeater
	TypeGenerator
)
