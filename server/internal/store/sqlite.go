package store

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	migracoes "github.com/dev-macb/mira/server/migrations"
	_ "modernc.org/sqlite"
)

func AbrirSQLite(contexto context.Context, caminhoBanco string) (*sql.DB, error) {
	caminhoAbsoluto, erro := filepath.Abs(caminhoBanco)
	if erro != nil {
		return nil, erro
	}
	if erro := os.MkdirAll(filepath.Dir(caminhoAbsoluto), 0700); erro != nil {
		return nil, erro
	}

	caminhoNaURI := filepath.ToSlash(caminhoAbsoluto)
	if !strings.HasPrefix(caminhoNaURI, "/") {
		caminhoNaURI = "/" + caminhoNaURI
	}
	parametrosSQLite := url.Values{"_pragma": {"busy_timeout(5000)", "foreign_keys(ON)"}}
	uriConexao := url.URL{Scheme: "file", Path: caminhoNaURI, RawQuery: parametrosSQLite.Encode()}
	bancoDeDados, erro := sql.Open("sqlite", uriConexao.String())
	if erro != nil {
		return nil, erro
	}
	// As operações da aplicação usam um só arquivo SQLite e um único escritor.
	bancoDeDados.SetMaxOpenConns(1)
	if erro := migracoes.Aplicar(contexto, bancoDeDados); erro != nil {
		bancoDeDados.Close()
		return nil, fmt.Errorf("aplicar migrações: %w", erro)
	}
	return bancoDeDados, nil
}
