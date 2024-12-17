package poi

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"hash/fnv"

	"google.golang.org/protobuf/proto"
)

var poi = []string{
	"中关村",
	"五道口",
	"知春路",
	"上地",
	"西二旗",
	"望京",
}

// Manager defines a poi manager.
type Manager struct {
}

// Resolve resolves the given location.
func (m *Manager) Resolve(c context.Context, loc *rentalpb.Location) (string, error) {
	b, err := proto.Marshal(loc)
	if err != nil {
		return "", err
	}
	h := fnv.New32()
	h.Write(b)
	return poi[int(h.Sum32())%len(poi)], nil
}
