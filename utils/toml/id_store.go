package toml

import "fmt"

type IdStore struct {
	idChain []string
	data    *Helper
}

func NewIdStore(d *Helper, idChain []string) *IdStore {
	return &IdStore{
		idChain: idChain,
		data:    d,
	}
}

func (p *IdStore) GetId() string {
	v, err := p.data.GetParamString("Version", "")
	if err != nil {
		return ""
	}

	return p.human2id(v)
}

func (p *IdStore) SetId(id string) {
	v := p.id2human(id)
	_, ok := p.data.GetTable().Fields["Version"]
	var err error
	if !ok {
		err = p.data.AddParam("Version", "", v)
	} else {
		err = p.data.SetParam("Version", "", v)
	}
	if err != nil {
		panic(err)
	}
}

func (p *IdStore) id2human(id string) string {
	for i, x := range p.idChain {
		if x != id {
			continue
		}
		return fmt.Sprintf("v.%d.0", i+1)
	}
	panic("id is not from idChain")
}

func (p *IdStore) human2id(str string) string {
	var i int
	_, err := fmt.Sscanf(str, "v.%d.0", &i)
	if err == nil && i > 0 && i <= len(p.idChain) {
		return p.idChain[i-1]
	}

	return str
}
