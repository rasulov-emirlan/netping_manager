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
	SNMPmib string `json:"snmpMib"`
}
