package cmd

// Actions enumeration
const (
	NOTHING = iota
	UPGRADE
	INSTALL
	REMOVE
	ERROR
	// Leaf bit
	LEAF     = 1 << 4
	ACT_MASK = LEAF - 1
)

type Action struct {
	Code    int
	Message string
}
