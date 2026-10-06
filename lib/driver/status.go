package driver

type DriverStatus int

const (
	Available DriverStatus = iota
	Busy
	Offline
)

var statusNames = map[DriverStatus]string{
	Available: "available",
	Busy:      "busy",
	Offline:   "offline",
}

func (s DriverStatus) String() string {
	if name, ok := statusNames[s]; ok {
		return name
	}
	return statusNames[Available]
}
