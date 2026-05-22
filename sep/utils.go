package sep

import (
	"fmt"
	"strconv"
	"time"
)

const (
	PollRate = 900
)

func GetTime() TimeType {
	return TimeType(time.Now().Second())
}

func ToSFDI(lfdi string) (SFDIType, error) {
	if len(lfdi) < 9 {
		return 0, fmt.Errorf("lfdi too short: %s", lfdi)
	}
	sfdi, err := strconv.ParseUint(lfdi[:9], 16, 64)
	if err != nil {
		return 0, err
	}

	return sfdi, nil
}

func NewLink(href string) *Link {
	return &Link{
		HrefAttr: href,
	}
}

func NewListLink(link *Link, all UInt32) *ListLink {
	return &ListLink{
		AllAttr: all,
		Link:    link,
	}
}

func NewResource(href string) *Resource {
	return &Resource{
		HrefAttr: href,
	}
}

func NewList(resource *Resource, all UInt32, result UInt32) *List {
	return &List{
		Resource:    resource,
		AllAttr:     all,
		ResultsAttr: result,
	}
}

func NewFunctionSetAssignmentsBase(resource *Resource) *FunctionSetAssignmentsBase {
	return &FunctionSetAssignmentsBase{
		Resource: resource,
	}
}

func NewDeviceCapability(fsab *FunctionSetAssignmentsBase) *DeviceCapability {
	return &DeviceCapability{
		PollRateAttr:               PollRate,
		FunctionSetAssignmentsBase: fsab,
	}
}

func NewAbstractDevice(resource *Resource, sfdi SFDIType) *AbstractDevice {
	return &AbstractDevice{
		SubscribableResource: &SubscribableResource{
			Resource: resource,
		},
		SFDI: &sfdi,
	}
}

func NewExternalDevice(adev *AbstractDevice) *ExternalDevice {
	t := time.Now().Unix()
	return &ExternalDevice{
		ChangedTime:    &t,
		Enabled:        true,
		AbstractDevice: adev,
	}
}

func NewEndDevice(ext *ExternalDevice) *EndDevice {
	return &EndDevice{
		ExternalDevice: ext,
	}
}

func NewEndDeviceList() *EndDeviceList {
	return &EndDeviceList{
		PollRateAttr:     PollRate,
		EndDevice:        make([]*EndDevice, 0),
		SubscribableList: &SubscribableList{},
	}
}

func NewRegistration(resource *Resource, dt *TimeType, pin *PINType) *Registration {
	return &Registration{
		Resource:           resource,
		DateTimeRegistered: dt,
		PIN:                pin,
		PollRateAttr:       PollRate,
	}
}

func NewTime(resource *Resource) *Time {
	t := time.Now()
	_, tz_offset := t.Zone()
	start, end := t.ZoneBounds()
	offset := 60 * 60

	// need to flip values and offset end by a year if not DST
	if !t.IsDST() {
		offset = 0
		start, end = end, start
		end.AddDate(1, 0, 0)
	}

	return &Time{
		PollRateAttr: PollRate,
		Resource:     resource,
		CurrentTime:  t.UTC().Unix(),
		DstStartTime: start.Unix(),
		DstEndTime:   end.Unix(),
		DstOffset:    offset,
		LocalTime:    t.Unix(),
		TzOffset:     tz_offset,
		Quality:      3,
	}
}

func NewFlowReservationRequestList(list *List) *FlowReservationRequestList {
	return &FlowReservationRequestList{
		List:                   list,
		PollRateAttr:           PollRate,
		FlowReservationRequest: make([]*FlowReservationRequest, 0),
	}
}
