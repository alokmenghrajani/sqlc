func myScan(out interface{}, rows *sql.Rows) error {
  columns, err := rows.Columns()
  if err != nil {
    return err
  }

  var values []reflect.Value // TODO: we know the size of values ahead of time.
  st := reflect.ValueOf(out).Elem().Type()
  for _, column := range columns {
    found := false
    for i := 0; i < st.NumField(); i++ {
      st_name := st.Field(i).Tag.Get("µorm")
      if strings.EqualFold(st_name, column) {
        values = append(values, reflect.ValueOf(out).Elem().Field(i).Addr())
        found = true
        break
      }
    }
    if !found {
      return errors.New(fmt.Sprintf("couldn't map %s", column))
    }
  }

  method := reflect.ValueOf(rows).MethodByName("Scan")
  result := method.Call(values)
  if !result[0].IsNil() {
    return result[0].Interface().(error)
  }
  return nil
}
