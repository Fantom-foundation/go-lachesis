package app

type (
	// StoreConfig is a config for store db.
	StoreConfig struct {
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

// DefaultStoreConfig for product.
func DefaultStoreConfig() StoreConfig {
	return StoreConfig{
		BlockCacheSize:      100,
		ReceiptsCacheSize:   100,
		DelegatorsCacheSize: 4000,
		StakersCacheSize:    4000,
	}
}

// LiteStoreConfig is for tests or inmemory.
func LiteStoreConfig() StoreConfig {
	return StoreConfig{
		BlockCacheSize:      50,
		ReceiptsCacheSize:   100,
		DelegatorsCacheSize: 400,
		StakersCacheSize:    400,
	}
}
