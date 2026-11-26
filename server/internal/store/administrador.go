package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func CriarAdministradorInicial(contexto context.Context, banco *sql.DB, nome, email, senha string) error {
	var quantidade int
	if erro := banco.QueryRowContext(contexto, "SELECT COUNT(*) FROM users").Scan(&quantidade); erro != nil {
		return erro
	}
	if quantidade > 0 {
		return nil
	}
	nome = strings.TrimSpace(nome)
	if nome == "" {
		nome = "Administrador"
	}
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || senha == "" {
		return errors.New("MIRA_ADMIN_EMAIL e MIRA_ADMIN_PASSWORD são obrigatórios enquanto não há usuários cadastrados")
	}
	endereco, erro := mail.ParseAddress(email)
	if erro != nil || endereco.Address != email || len(email) > 254 {
		return errors.New("MIRA_ADMIN_EMAIL inválido")
	}
	if len(nome) > 100 || len(senha) < 12 || len(senha) > 72 {
		return errors.New("nome deve ter até 100 caracteres e a senha inicial entre 12 e 72 bytes")
	}
	resumoSenha, erro := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	if erro != nil {
		return fmt.Errorf("gerar hash de senha: %w", erro)
	}
	_, erro = banco.ExecContext(contexto, "INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?)", nome, email, string(resumoSenha))
	if erro != nil {
		return fmt.Errorf("gravar administrador inicial: %w", erro)
	}
	return nil
}
