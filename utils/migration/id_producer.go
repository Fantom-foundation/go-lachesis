package migration

type IdProducer interface {
	GetId() string
	SetId(string)
}
