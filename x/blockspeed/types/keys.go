package types

const (
	// ModuleName defines the module name
	ModuleName = "blockspeed"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_blockspeed"
)

var (
	ParamsKey = []byte("p_blockspeed")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
