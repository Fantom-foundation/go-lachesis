package version

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/params"
)

func init() {
	params.VersionMajor = 1     // Major version component of the current release
	params.VersionMinor = 0     // Minor version component of the current release
	params.VersionPatch = 0     // Patch version component of the current release
	params.VersionMeta = "rc.0" // Version metadata to append to the version string
}

func BigToString(b *big.Int) string {
	if len(b.Bytes()) > 8 {
		return "_malformed_version_"
	}
	return U64ToString(b.Uint64())
}

func U64ToString(v uint64) string {
	return fmt.Sprintf("%d.%d.%d", (v/1e12)%1e6, (v/1e6)%1e6, v%1e6)
}

func AsString() string {
	return fmt.Sprintf("%d.%d.%d", params.VersionMajor, params.VersionMinor, params.VersionPatch)
}

func AsU64() uint64 {
	return asU64(uint16(params.VersionMajor), uint16(params.VersionMinor), uint16(params.VersionPatch))
}

func AsBigInt() *big.Int {
	return new(big.Int).SetUint64(AsU64())
}

func asU64(vMajor, vMinor, vPatch uint16) uint64 {
	return uint64(vMajor)*1e12 + uint64(vMinor)*1e6 + uint64(vPatch)
}
