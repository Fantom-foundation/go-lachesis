package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/naoina/toml/ast"

	"github.com/Fantom-foundation/go-lachesis/utils/migration"
)

/*
	Methods for migrations (*ConfigData):
	- AddSection
	- DeleteSection
	- RenameSection

	- AddParam
	- DeleteParam
	- RenameParam
	- SetParam
	- GetParam[String|Int|Float|Bool|Time]
*/

func ConfigMigrations(data *ConfigData) *migration.Migration {
	return migration.Begin("lachesis-config")

	/*
		Use here only named migrations. Migration name - version of config.
		Example ():

		  return migration.Init("lachesis-config", "ajIr@Quicuj9"
			).NewNamed("v1", func()error{
				... // Some actions for migrations
				return err
			}).NewNamed("v2", func()error{
				... // Some actions for migrations
				return err
			})
			...
	*/

}

type ConfigData struct {
	table *ast.Table
}

func NewConfigData(t *ast.Table) *ConfigData {
	return &ConfigData{
		table: t,
	}
}

func (d *ConfigData) GetTable() *ast.Table {
	return d.table
}

func (d *ConfigData) AddSection(name, after string) error {
	if name == "" {
		return nil
	}

	_, err := d.findSection(name)
	if err == nil {
		// If exists - return error
		return errors.New("section already exists: " + name)
	}

	path := strings.Split(name, ".")

	afterSection, err := d.findSection(after)
	if err != nil {
		return err
	}

	pathStr := ""
	currentSection := d.table
	for _, n := range path {
		pathStr = pathStr + "/" + n

		var section *ast.Table
		sectionI, ok := currentSection.Fields[n]
		if ok {
			section, ok = sectionI.(*ast.Table)
			if !ok {
				return errors.New("wrong type of section: " + pathStr)
			}
		} else {
			section = &ast.Table{
				Position: ast.Position{
					Begin: afterSection.End() + 1,
					End:   0,
				},
				Line:   afterSection.Line + afterSection.End() - afterSection.Pos(),
				Name:   n,
				Fields: make(map[string]interface{}),
				Type:   ast.TableTypeNormal,
			}

		}
		currentSection.Fields[n] = section
		currentSection = section
	}

	return nil
}

func (d *ConfigData) DeleteSection(name string) error {
	// Find parent section and name for deletedName section
	path := strings.Split(name, ".")
	parentName := strings.Join(path[:len(path)-1], ".")
	deletedName := path[len(path)-1]

	parent, err := d.findSection(parentName)
	if err != nil {
		return err
	}

	delete(parent.Fields, deletedName)

	return nil
}

func (d *ConfigData) RenameSection(name, newName string) error {
	section, err := d.findSection(name)
	if err != nil {
		return err
	}

	section.Name = newName
	delete(d.table.Fields, name)
	d.table.Fields[newName] = section

	return nil
}

func (d *ConfigData) AddParam(name, sectionName string, value interface{}) error {
	_, sect, err := d.getKVData(name, sectionName)
	if err == nil {
		return errors.New("param already exists in section: " + sectionName + " / " + name)
	}
	if sect == nil {
		return err
	}

	kvData, err := d.setKVData(name, value)
	if err != nil {
		return err
	}

	sect.Fields[name] = kvData

	return nil
}

func (d *ConfigData) DeleteParam(name, sectionName string) error {
	_, sect, err := d.getKVData(name, sectionName)
	if err != nil {
		return err
	}

	delete(sect.Fields, name)

	return nil
}

func (d *ConfigData) RenameParam(name, sectionName, newName string) error {
	param, sect, err := d.getKVData(name, sectionName)
	if err != nil {
		return err
	}

	delete(sect.Fields, name)
	param.Key = newName
	sect.Fields[newName] = param

	return nil
}

func (d *ConfigData) SetParam(name, sectionName string, value interface{}) error {
	param, _, err := d.getKVData(name, sectionName)
	if err != nil {
		return err
	}

	_, err = d.setKVData(name, value, param)
	if err != nil {
		return err
	}

	return nil
}

func (d *ConfigData) GetParamString(name, sectionName string) (string, error) {
	param, _, err := d.getKVData(name, sectionName)
	if err != nil {
		return "", err
	}

	pString, ok := param.Value.(*ast.String)
	if !ok {
		return "", errors.New("wrong type for string in param: " + sectionName + " / " + name)
	}

	return pString.Value, nil
}

func (d *ConfigData) GetParamInt(name, sectionName string) (int64, error) {
	param, _, err := d.getKVData(name, sectionName)
	if err != nil {
		return -1, err
	}

	pInt, ok := param.Value.(*ast.Integer)
	if !ok {
		return -1, errors.New("wrong type for integer in param: " + sectionName + " / " + name)
	}

	return pInt.Int()
}

func (d *ConfigData) GetParamFloat(name, sectionName string) (float64, error) {
	param, _, err := d.getKVData(name, sectionName)
	if err != nil {
		return -1, err
	}

	pFloat, ok := param.Value.(*ast.Float)
	if !ok {
		return -1, errors.New("wrong type for integer in param: " + sectionName + " / " + name)
	}

	return pFloat.Float()
}

func (d *ConfigData) GetParamBool(name, sectionName string) (bool, error) {
	param, _, err := d.getKVData(name, sectionName)
	if err != nil {
		return false, err
	}

	pBool, ok := param.Value.(*ast.Boolean)
	if !ok {
		return false, errors.New("wrong type for integer in param: " + sectionName + " / " + name)
	}

	return pBool.Boolean()
}

func (d *ConfigData) GetParamTime(name, sectionName string) (time.Time, error) {
	param, _, err := d.getKVData(name, sectionName)
	if err != nil {
		return time.Now(), err
	}

	pTime, ok := param.Value.(*ast.Datetime)
	if !ok {
		return time.Now(), errors.New("wrong type for integer in param: " + sectionName + " / " + name)
	}

	return pTime.Time()
}

func (d *ConfigData) findSection(name string) (*ast.Table, error) {
	path := strings.Split(name, ".")
	currentSection := d.table

	if name == "" {
		return currentSection, nil
	}

	pathStr := ""
	for _, n := range path {
		pathStr = pathStr + "/" + n

		sectionI, ok := currentSection.Fields[n]
		if !ok {
			return nil, errors.New("section not found: " + pathStr)
		}

		currentSection, ok = sectionI.(*ast.Table)
		if !ok {
			return nil, errors.New("section has wrong type: " + pathStr)
		}
	}

	return currentSection, nil
}

func (d *ConfigData) getKVData(name, sectionName string) (*ast.KeyValue, *ast.Table, error) {
	sect, err := d.findSection(sectionName)
	if err != nil {
		return nil, nil, err
	}
	if sect == nil {
		return nil, nil, errors.New("section not found: " + sectionName)
	}

	paramI, ok := sect.Fields[name]
	if !ok {
		return nil, sect, errors.New("param not exists in section: " + sectionName + " / " + name)
	}

	param, ok := paramI.(*ast.KeyValue)
	if !ok {
		return nil, sect, errors.New("wrong param type in section: " + sectionName + " / " + name)
	}

	return param, sect, nil
}

func (d *ConfigData) setKVData(name string, value interface{}, kvExists ...*ast.KeyValue) (*ast.KeyValue, error) {
	var kv *ast.KeyValue
	if len(kvExists) > 0 {
		kv = kvExists[0]
	} else {
		kv = &ast.KeyValue{
			Key: name,
		}
	}
	switch value.(type) {
	case string:
		kv.Value = &ast.String{
			Position: ast.Position{},
			Value:    value.(string),
			Data:     []rune(value.(string)),
		}
	case int:
		s := strconv.FormatInt(int64(value.(int)), 10)
		kv.Value = &ast.Integer{
			Position: ast.Position{},
			Value:    s,
			Data:     []rune(s),
		}
	case float64:
		s := strconv.FormatFloat(value.(float64), 'f', 16, 64)
		kv.Value = &ast.Float{
			Position: ast.Position{},
			Value:    s,
			Data:     []rune(s),
		}
	case bool:
		s := strconv.FormatBool(value.(bool))
		kv.Value = &ast.Boolean{
			Position: ast.Position{},
			Value:    s,
			Data:     []rune(s),
		}
	case time.Time:
		s := value.(time.Time).Format("2006-01-02T15:04:05.999999999Z07:00")
		kv.Value = &ast.Datetime{
			Position: ast.Position{},
			Value:    s,
			Data:     []rune(s),
		}
	}

	return kv, nil
}

type tomlIdStore struct {
	idChain []string
	data    *ConfigData
}

func NewTomlIdStore(d *ConfigData, idChain []string) *tomlIdStore {
	return &tomlIdStore{
		idChain: idChain,
		data:    d,
	}
}

func (p *tomlIdStore) GetId() string {
	v, err := p.data.GetParamString("Version", "")
	if err != nil {
		return ""
	}

	return p.human2id(v)
}

func (p *tomlIdStore) SetId(id string) {
	v := p.id2human(id)
	_, ok := p.data.GetTable().Fields["Version"]
	if !ok {
		err := p.data.AddParam("Version", "", v)
		if err != nil {
			panic(err)
		}
	} else {
		err := p.data.SetParam("Version", "", v)
		if err != nil {
			panic(err)
		}
	}
}

func (p *tomlIdStore) id2human(id string) string {
	for i, x := range p.idChain {
		if x != id {
			continue
		}
		return fmt.Sprintf("v.%d.0", i+1)
	}
	panic("id2human() fail")
}

func (p *tomlIdStore) human2id(str string) string {
	var i int
	_, err := fmt.Sscanf(str, "v.%d.0", &i)
	if err != nil {
		panic(err)
	}
	return p.idChain[i-1]
}
