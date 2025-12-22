package main

import (
	"log"
	"reflect"
	"strings"

	"cloud.google.com/go/civil"
	"github.com/go-jet/jet/v2/generator/metadata"
	"github.com/go-jet/jet/v2/generator/postgres"
	"github.com/go-jet/jet/v2/generator/template"
	postgres2 "github.com/go-jet/jet/v2/postgres"
	_ "github.com/lib/pq"
)

var skipTableName = map[string]struct{}{
	"schema_migrations": {},
}

// Have to make this nicer
// https://github.com/go-jet/jet/wiki/Generator#generator-customization
func main() {
	err := postgres.Generate(
		"./internal/sqlgen",
		postgres.DBConnection{
			Host:       "localhost",
			Port:       5430,
			User:       "postgres",
			Password:   "postgrespw",
			DBName:     "bigbackend",
			SchemaName: "public",
			SslMode:    "disable",
		},
		template.Default(postgres2.Dialect).
			UseSchema(func(schema metadata.Schema) template.Schema {
				return template.DefaultSchema(schema).
					// Generate SQL Builders
					UseSQLBuilder(template.DefaultSQLBuilder().
						UseTable(func(table metadata.Table) template.TableSQLBuilder {

							_, ok := skipTableName[table.Name]
							if ok {
								return template.TableSQLBuilder{
									Skip: true,
								}
							}

							return template.DefaultTableSQLBuilder(table)
						}),
					).
					// Generate models
					UseModel(template.DefaultModel().
						UsePath("/entity").
						UseTable(func(table metadata.Table) template.TableModel {

							_, ok := skipTableName[table.Name]
							if ok {
								return template.TableModel{
									Skip: true,
								}
							}

							return template.DefaultTableModel(table).
								UseField(func(column metadata.Column) template.TableModelField {
									defaultTableModelField := template.DefaultTableModelField(column)

									goType, found := getGoTypeOverride(column)
									if found {
										defaultTableModelField.Type = template.NewType(goType)
										return defaultTableModelField
									}

									return defaultTableModelField
								})
						}),
					)
			}),
	)

	if err != nil {
		log.Println("failed to generate models:", err)
		return
	}
	return
}

func getGoTypeOverride(column metadata.Column) (any, bool) {
	defaultGoType, found := toGoTypeOverride(column)
	if !found {
		return nil, false
	}

	if column.IsNullable {
		return reflect.New(reflect.TypeOf(defaultGoType)).Interface(), true
	}

	return defaultGoType, true
}

func toGoTypeOverride(column metadata.Column) (any, bool) {
	switch strings.ToLower(column.Name) {
	case "date_of_birth":
		return civil.Date{}, true
	default:
		return nil, false
	}
}
