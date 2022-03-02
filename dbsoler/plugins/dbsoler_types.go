package main

import (
	"reflect"

	"github.com/lapuglisi/dbsoler/dbsol"
	"github.com/lapuglisi/dbsoler/sra"
)

// InitTypes .....
func InitTypes(types map[string]reflect.Type) {
	// dbsoler.dbsol
	types["variaveis"] = reflect.TypeOf(dbsol.Variaveis{})
	types["entidades"] = reflect.TypeOf(dbsol.Entidades{})
	types["transacoes"] = reflect.TypeOf(dbsol.Transacoes{})
	types["tm_tabelas"] = reflect.TypeOf(dbsol.Tabelas{})
	types["processos"] = reflect.TypeOf(dbsol.Processos{})
	types["formas_comunicacao"] = reflect.TypeOf(dbsol.FormasComunicao{})

	// dbsoler.sra
	types["sub_prod_adq_ec"] = reflect.TypeOf(sra.SubProdAdqEc{})
}
