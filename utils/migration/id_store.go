package migration

type IdStore interface {
	GetId() string
	SetId(string)
}

type inmemIdStore struct {
	lastId string
}

func (p *inmemIdStore) GetId() string {
	return string(p.lastId)
}

func (p *inmemIdStore) SetId(id string) {
	p.lastId = id
}
