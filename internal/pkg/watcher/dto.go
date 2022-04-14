package watcher

import "github.com/rasulov-emirlan/netping-manager/internal/manager"

func locationsToService(l []*Location) []*manager.Location {
	res := make([]*manager.Location, len(l))
	for i, v := range l {
		res[i] = v.model
	}
	return res
}
