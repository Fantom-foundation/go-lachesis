package migration

type IdProducer interface {
	GetId() string
	SetId(string)
	IsCurrent(string) bool
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

func (p *inmemIdProducer) IsCurrent(id string) bool {
	currentId := p.GetId()
	return id == currentId
}
