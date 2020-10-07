package version

import (
	"math/big"

	"github.com/ethereum/go-ethereum/params"
)

func init() {
	params.VersionMajor = 0     // Major version component of the current release
	params.VersionMinor = 7     // Minor version component of the current release
	params.VersionPatch = 0     // Patch version component of the current release
	params.VersionMeta = "rc.1" // Version metadata to append to the version string
}

func AsBigInt() *big.Int {
	return asBigInt(uint64(params.VersionMajor), uint64(params.VersionMinor), uint64(params.VersionPatch))
}

func asBigInt(vMajor, vMinor, vPatch uint64) *big.Int {
	return new(big.Int).Add(
		new(big.Int).Add(
			new(big.Int).Lsh(
				new(big.Int).SetUint64(vMajor), 64*3),
			new(big.Int).Lsh(
				new(big.Int).SetUint64(vMinor), 64*2),
		), new(big.Int).Lsh(
			new(big.Int).SetUint64(vPatch), 64*1),
		// VersionMeta is not used here
	)
}
