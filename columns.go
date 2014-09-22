// THIS FILE WAS AUTOGENERATED - ANY EDITS TO THIS WILL BE LOST WHEN IT IS REGENERATED

package sqlc



type varcharColumn struct {
	name string
}

type VarcharColumn interface {
	Column
	Eq(value string) Condition
}

func (c *varcharColumn) ColumnName() string {
	return c.name
}

func (c *varcharColumn) Eq(pred string) Condition {
	return Condition{Binding: ColumnBinding{Value: pred, Column: c}}
}

func VarcharField(name string) VarcharColumn {
	return &varcharColumn{name: name}
}

