package app

type (
	// Config for the application.
	Config struct {
		StoreConfig
	}

	// StoreConfig is a config for store db.
	StoreConfig struct {
		// TxIndex enables indexing transactions and receipts.
		TxIndex bool
		// Cache size for Block.
		BlockCacheSize int
		// Cache size for Receipts.
		ReceiptsCacheSize int
		// Cache size for Stakers.
		StakersCacheSize int
		// Cache size for Delegators.
		DelegatorsCacheSize int
	}
)

// DefaultConfig for product.
func DefaultConfig() Config {
	return Config{
		DefaultStoreConfig(),
	}
}

// DefaultStoreConfig for product.
func DefaultStoreConfig() StoreConfig {
	return StoreConfig{
		TxIndex:             true,
		BlockCacheSize:      100,
		ReceiptsCacheSize:   100,
		DelegatorsCacheSize: 4000,
		StakersCacheSize:    4000,
	}
}

// LiteStoreConfig is for tests or inmemory.
func LiteStoreConfig() StoreConfig {
	return StoreConfig{
		TxIndex:             true,
		BlockCacheSize:      50,
		ReceiptsCacheSize:   100,
		DelegatorsCacheSize: 400,
		StakersCacheSize:    400,
	}
}
