package memory

import (
	"github.com/Tylores/egot/internal/sep"
	"github.com/Tylores/egot/internal/uri"
)

const (
	PollRate = 10
)

var (
	StaticDeviceCapability = sep.DeviceCapability{
		PollRateAttr: PollRate,
		FunctionSetAssignmentsBase: &sep.FunctionSetAssignmentsBase{
			Resource: &sep.Resource{
				HrefAttr: uri.DeviceCapability,
			},
			TimeLink: &sep.TimeLink{
				Link: &sep.Link{
					HrefAttr: uri.Time,
				},
			},
		},
		EndDeviceListLink: &sep.EndDeviceListLink{
			ListLink: &sep.ListLink{
				Link: &sep.Link{
					HrefAttr: uri.EndDeviceList,
				},
				AllAttr: 0,
			},
		},
	}
)
