package migrations

import (
	"context"
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var arquivosMigracao embed.FS

func Aplicar(contexto context.Context, bancoDeDados *sql.DB) error {
	provedorMigracoes, erro := goose.NewProvider(goose.DialectSQLite3, bancoDeDados, arquivosMigracao)
	if erro != nil {
		return erro
	}
	_, erro = provedorMigracoes.Up(contexto)
	return erro
}
