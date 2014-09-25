package sqlc

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

var predicateTypes = map[PredicateType]string{
	EqPredicate: "=",
	GtPredicate: ">",
	GePredicate: ">=",
	LtPredicate: "<",
	LePredicate: "<=",
}

func (s *selection) String() string {
	var buf bytes.Buffer
	s.Render(&buf)
	return buf.String()
}

func (s *selection) Render(w io.Writer) (placeholders []interface{}) {
	fmt.Fprint(w, "SELECT ")

	if len(s.projection) == 0 {
		fmt.Fprint(w, "*")
	} else {
		colClause := columnClause(s.projection)
		fmt.Fprint(w, colClause)
	}

	fmt.Fprintf(w, " FROM ")

	switch sub := s.selection.(type) {
	case table:
		fmt.Fprint(w, sub.name)
	case *selection:
		fmt.Fprint(w, "(")
		sub.Render(w)
		fmt.Fprint(w, ")")
	}

	if len(s.predicate) > 0 {
		fmt.Fprint(w, " ")
		placeholders = renderWhereClause(s.predicate, w)
	} else {
		placeholders = []interface{}{}
	}

	if (len(s.groups)) > 0 {
		fmt.Fprint(w, " GROUP BY ")
		colClause := columnClause(s.groups)
		fmt.Fprint(w, colClause)
	}

	// TODO eliminate copy and paste
	if (len(s.ordering)) > 0 {
		fmt.Fprint(w, " ORDER BY ")
		colClause := columnClause(s.ordering)
		fmt.Fprint(w, colClause)
	}

	return placeholders
}

func columnClause(cols []Field) string {
	colFragments := make([]string, len(cols))
	for i, col := range cols {
		colFragments[i] = col.Name()
	}
	return strings.Join(colFragments, ", ")
}

func renderWhereClause(conds []Condition, w io.Writer) []interface{} {
	fmt.Fprint(w, "WHERE ")

	whereFragments := make([]string, len(conds))
	values := make([]interface{}, len(conds))

	for i, condition := range conds {
		col := condition.Binding.Field.Name()
		pred := condition.Predicate
		whereFragments[i] = fmt.Sprintf("%s %s ?", col, predicateTypes[pred])
		values[i] = condition.Binding.Value
	}

	whereClause := strings.Join(whereFragments, " AND ")
	fmt.Fprint(w, whereClause)

	return values
}
