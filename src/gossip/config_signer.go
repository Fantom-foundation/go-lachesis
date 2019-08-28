package gossip

type SignerConfig struct {
	KeystoreLocation string // location of keystore directory
}

func DefaultSignerConfig() SignerConfig {
	return SignerConfig{
		KeystoreLocation: "lachesis-config",
	}
}
