package fields

type ProtocolVersion string

const (
	ProtocolVersion1 ProtocolVersion = "1.00"
	ProtocolVersion2 ProtocolVersion = "2.00"
)

func (p ProtocolVersion) ID() string {
	return string(p)
}

func (p ProtocolVersion) String() string {
	return string(p)
}
