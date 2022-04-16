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

const (
	TypeAC = iota + 1
	TypeHeater
	TypeGenerator
)
