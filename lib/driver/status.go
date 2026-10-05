package driver

type DriverStatus int

const (
	Available DriverStatus = iota
	Busy
	Offline
)

func (s DriverStatus) String() string {
	switch s {
	case Available:
		return "available"
	case Busy:
		return "busy"
	case Offline:
		return "offline"
	default:
		return "available"
	}
}
