package migration

type IdProducer interface {
	GetId() (string, error)
	SetId(string) error
}
