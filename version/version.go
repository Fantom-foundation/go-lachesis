package version

import (
	"math/big"

	"github.com/ethereum/go-ethereum/params"
)

func init() {
	params.VersionMajor = 0     // Major version component of the current release
	params.VersionMinor = 8     // Minor version component of the current release
	params.VersionPatch = 0     // Patch version component of the current release
	params.VersionMeta = "rc.1" // Version metadata to append to the version string
}

func AsU64() uint64 {
	return asU64(uint16(params.VersionMajor), uint16(params.VersionMinor), uint16(params.VersionPatch))
}

func AsBigInt() *big.Int {
	return new(big.Int).SetUint64(AsU64())
}

func asU64(vMajor, vMinor, vPatch uint16) uint64 {
	return uint64(vMajor) * 1e12 + uint64(vMinor) * 1e6 + uint64(vPatch)
}
