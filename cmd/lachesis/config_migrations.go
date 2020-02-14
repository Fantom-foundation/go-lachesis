package main

import (
	"github.com/Fantom-foundation/go-lachesis/utils/migration"
	"github.com/naoina/toml/ast"
)

func addTablePart(t *ast.Table, partName string) *ast.Table {
	newPart := &ast.Table{
		Name:     partName,
		Fields:   make(map[string]interface{}),
		Type:     ast.TableTypeNormal,
		Data:     nil,
	}

	t.Fields[partName] = newPart

	return newPart
}

func addTableKeyString(t *ast.Table, key string) *ast.KeyValue {
	newKV := &ast.KeyValue{
		Key:   key,
		Value: &ast.String{
			Value:    "",
		},
		Line:  0,
	}

	t.Fields[key] = newKV

	return newKV
}

func ConfigMigrations(mapTable *ast.Table) *migration.Migration {
	return migration.Init("lachesis-config", "ajIr@Quicuj9")

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

type configIdProducer struct {
	table *ast.Table
}

func NewConfigIdProducer(t *ast.Table) *configIdProducer {
	return &configIdProducer{
		table: t,
	}
}

func (p *configIdProducer) GetId() (string, error) {
	versionI, ok := p.table.Fields["Version"]
	if !ok {
		return "", nil
	}
	versionStr := versionI.(*ast.KeyValue).Value.(*ast.String).Value

	return versionStr, nil
}

func (p *configIdProducer) SetId(id string) error {
	// fmt.Printf("DBG:\n%+v\n", p.table.Fields["Lachesis"].(*ast.Table).Fields["EVMInterpreter"].(*ast.KeyValue).Value.(*ast.String))
	configVersionI, ok := p.table.Fields["Version"]
	if !ok {
		configVersionI = addTableKeyString(p.table, "Version")
	}
	configVersion := configVersionI.(*ast.KeyValue)
	configVersionValue, ok := configVersion.Value.(*ast.String)
	if !ok {
		panic("Bad type field Version on config. Must by String!")
	}
	configVersionValue.Value = id

	// fmt.Printf("DBG:\n%+v\n", p.table)

	return nil
}
