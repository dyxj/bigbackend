package main

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"cloud.google.com/go/civil"
	"github.com/dyxj/bigbackend/pkg/sqldb"
	"github.com/go-jet/jet/v2/generator/metadata"
	"github.com/go-jet/jet/v2/generator/postgres"
	"github.com/go-jet/jet/v2/generator/template"
	postgres2 "github.com/go-jet/jet/v2/postgres"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
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
										defaultTableModelField.Type = buildTemplateType(goType)
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

func buildTemplateType(v any) template.Type {
	t := reflect.TypeOf(v)
	if !isGeneric(t) {
		return template.NewType(v)
	}
	return buildGenericTemplateType(t)
}

func isGeneric(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// 1. Slices/Maps return "" for Name()
	// 2. Custom types like 'type MySlice []int' return "MySlice", no brackets.
	// 3. Only instantiated generics return "TypeName[...]"
	name := t.Name()
	return name != "" && strings.Contains(name, "[")
}

// Given A[x/y/z.B]
func buildGenericTemplateType(t reflect.Type) template.Type {
	isPointer := false
	for t.Kind() == reflect.Ptr {
		isPointer = true
		t = t.Elem()
	}

	name := t.String()
	parts := strings.Split(name, "[")

	// Obtain A
	genericType := parts[0]
	if isPointer {
		genericType = "*" + genericType
	}

	// Obtain x/y/z.B
	isPointerTypeArg := false
	typeArgWithImport := strings.TrimSuffix(parts[1], "]")
	if strings.Contains(typeArgWithImport, "*") {
		isPointerTypeArg = true
		typeArgWithImport = strings.TrimPrefix(typeArgWithImport, "*")
	}

	typeArgParts := strings.Split(typeArgWithImport, "/")

	// Obtain z.B
	typeArg := typeArgParts[len(typeArgParts)-1]

	// Obtain z
	typeArgPackageName := strings.Split(typeArg, ".")[0]

	if isPointerTypeArg {
		typeArg = "*" + typeArg
	}

	// Construct A[z.B]
	typeName := fmt.Sprintf("%s[%s]", genericType, typeArg)

	// Construct x/y/z
	typeArgImport := strings.Join(
		append(typeArgParts[:len(typeArgParts)-1], typeArgPackageName),
		"/",
	)

	return template.Type{
		ImportPath:            t.PkgPath(),
		AdditionalImportPaths: []string{typeArgImport},
		Name:                  typeName,
	}
}

func getGoTypeOverride(column metadata.Column) (any, bool) {
	defaultGoType, found := toGoTypeOverride(column)
	if !found {
		return nil, false
	}

	if column.DataType.IsArray() {
		defaultGoType = toGoArrayType(defaultGoType, column)
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
	}

	switch strings.ToLower(column.DataType.Name) {
	case "numeric":
		return decimal.Decimal{}, true
	}

	return nil, false
}

func toGoArrayType(elemType any, column metadata.Column) any {
	if column.DataType.Dimensions > 1 {
		return "" // unsupported multidimensional arrays
	}

	switch elemType.(type) {
	case bool:
		return pq.BoolArray{}
	case int32:
		return pq.Int32Array{}
	case int64:
		return pq.Int64Array{}
	case float32:
		return pq.Float32Array{}
	case float64:
		return pq.Float64Array{}
	case []byte:
		return pq.ByteaArray{}
	case decimal.Decimal:
		return sqldb.Array[decimal.Decimal]{}
	default:
		return pq.StringArray{}
	}
}
