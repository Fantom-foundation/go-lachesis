package migration

type IdProducer interface {
	GetId() string
	SetId(string)
}

type inmemIdProducer struct {
	lastId string
}

func (p *inmemIdProducer) GetId() string {
	return string(p.lastId)
}

func (p *inmemIdProducer) SetId(id string) {
	p.lastId = id
}
