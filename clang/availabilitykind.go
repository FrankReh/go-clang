package clang

// #include "go-clang.h"
import "C"

// Describes the availability of a particular entity, which indicates whether
// the use of this entity will result in a warning or error due to it being
// deprecated or unavailable.
type AvailabilityKind uint32

const (
	// The entity is available.
	Availability_Available AvailabilityKind = C.CXAvailability_Available

	// The entity is available, but has been deprecated (and its use is not recommended).
	Availability_Deprecated AvailabilityKind = C.CXAvailability_Deprecated

	// The entity is not available; any use of it will be an error.
	Availability_NotAvailable AvailabilityKind = C.CXAvailability_NotAvailable

	// The entity is available, but not accessible; any use of it will be an error.
	Availability_NotAccessible AvailabilityKind = C.CXAvailability_NotAccessible
)
