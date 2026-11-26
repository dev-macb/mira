package api

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	nomeCookieSessao      = "mira_sessao"
	duracaoTokenAcesso    = 15 * time.Minute
	duracaoTokenRenovacao = 30 * 24 * time.Hour
)

type contaAutenticada struct {
	ID           int64      `json:"id"`
	Nome         string     `json:"nome"`
	Email        string     `json:"email"`
	Perfil       string     `json:"perfil"`
	UltimoAcesso *time.Time `json:"ultimoAcesso"`
}

type limiteTentativas struct {
	quantidade int
	inicio     time.Time
}
type limitadorLogin struct {
	mu        sync.Mutex
	porOrigem map[string]limiteTentativas
	agora     func() time.Time
}

func (l *limitadorLogin) permitido(origem string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	agora := l.agora()
	registro := l.porOrigem[origem]
	if agora.Sub(registro.inicio) >= time.Minute {
		registro = limiteTentativas{inicio: agora}
	}
	for ip, tentativa := range l.porOrigem {
		if agora.Sub(tentativa.inicio) >= time.Minute {
			delete(l.porOrigem, ip)
		}
	}
	if registro.quantidade >= 10 {
		l.porOrigem[origem] = registro
		return false
	}
	registro.quantidade++
	l.porOrigem[origem] = registro
	return true
}

func novoRoteadorAutenticacao(banco *sql.DB, segredo string, agora func() time.Time) http.Handler {
	r := chi.NewRouter()
	limite := &limitadorLogin{porOrigem: make(map[string]limiteTentativas), agora: agora}
	r.Post("/login", func(w http.ResponseWriter, req *http.Request) {
		origem := req.Header.Get("X-Forwarded-For")
		if origem == "" {
			origem = req.RemoteAddr
		}
		if !limite.permitido(origem) {
			responderJSON(w, http.StatusTooManyRequests, map[string]string{"erro": "Muitas tentativas. Aguarde um minuto."})
			return
		}
		var entrada struct {
			Email string `json:"email"`
			Senha string `json:"senha"`
		}
		if json.NewDecoder(http.MaxBytesReader(w, req.Body, 4096)).Decode(&entrada) != nil {
			responderJSON(w, http.StatusBadRequest, map[string]string{"erro": "Informe e-mail e senha válidos."})
			return
		}
		var conta contaAutenticada
		var hash string
		var ativo bool
		var ultimo sql.NullTime
		err := banco.QueryRowContext(req.Context(), "SELECT id,name,email,role,password_hash,active,last_sign_in_at FROM users WHERE email=?", strings.ToLower(strings.TrimSpace(entrada.Email))).Scan(&conta.ID, &conta.Nome, &conta.Email, &conta.Perfil, &hash, &ativo, &ultimo)
		if err != nil || !ativo || bcrypt.CompareHashAndPassword([]byte(hash), []byte(entrada.Senha)) != nil {
			responderJSON(w, http.StatusUnauthorized, map[string]string{"erro": "E-mail ou senha incorretos."})
			return
		}
		instante := agora().UTC()
		if _, err = banco.ExecContext(req.Context(), "UPDATE users SET last_sign_in_at=? WHERE id=?", instante, conta.ID); err != nil {
			responderJSON(w, http.StatusInternalServerError, map[string]string{"erro": "Não foi possível iniciar a sessão."})
			return
		}
		conta.UltimoAcesso = &instante
		if err = emitirSessao(w, req, banco, segredo, conta, agora); err != nil {
			responderJSON(w, http.StatusInternalServerError, map[string]string{"erro": "Não foi possível iniciar a sessão."})
			return
		}
		responderJSON(w, http.StatusOK, conta)
	})
	autenticado := chi.NewRouter()
	autenticado.Use(middlewareSessao(banco, segredo, agora))
	autenticado.Get("/me", func(w http.ResponseWriter, req *http.Request) {
		responderJSON(w, http.StatusOK, req.Context().Value(chaveConta{}).(contaAutenticada))
	})
	autenticado.Post("/refresh", func(w http.ResponseWriter, req *http.Request) {
		conta := req.Context().Value(chaveConta{}).(contaAutenticada)
		if err := emitirSessao(w, req, banco, segredo, conta, agora); err != nil {
			responderJSON(w, http.StatusInternalServerError, map[string]string{"erro": "Não foi possível renovar a sessão."})
			return
		}
		responderJSON(w, http.StatusOK, conta)
	})
	autenticado.Post("/logout", func(w http.ResponseWriter, req *http.Request) {
		if c, e := req.Cookie(nomeCookieSessao); e == nil {
			h := sha256.Sum256([]byte(c.Value))
			_, _ = banco.ExecContext(req.Context(), "UPDATE refresh_tokens SET revoked=1 WHERE token_hash=?", base64.RawURLEncoding.EncodeToString(h[:]))
		}
		expirarCookieSessao(w, seguroRequisicao(req))
		responderJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	r.Mount("/", autenticado)
	return r
}

type chaveConta struct{}

func middlewareSessao(b *sql.DB, segredo string, agora func() time.Time) func(http.Handler) http.Handler {
	return func(proximo http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(nomeCookieSessao)
			if err != nil {
				responderJSON(w, http.StatusUnauthorized, map[string]string{"erro": "Sessão inválida."})
				return
			}
			hash := sha256.Sum256([]byte(cookie.Value))
			var id int64
			var nome, email, perfil string
			var ultimo sql.NullTime
			var ativo bool
			err = b.QueryRowContext(r.Context(), `SELECT u.id,u.name,u.email,u.role,u.last_sign_in_at,u.active FROM refresh_tokens t JOIN users u ON u.id=t.user_id WHERE t.token_hash=? AND t.revoked=0 AND t.expires_at>?`, base64.RawURLEncoding.EncodeToString(hash[:]), agora().UTC()).Scan(&id, &nome, &email, &perfil, &ultimo, &ativo)
			if err != nil || !ativo {
				expirarCookieSessao(w, seguroRequisicao(r))
				responderJSON(w, http.StatusUnauthorized, map[string]string{"erro": "Sessão inválida."})
				return
			}
			if !strings.HasSuffix(r.URL.Path, "/refresh") && !validarTokenAcesso(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "), segredo, id, agora()) {
				responderJSON(w, http.StatusUnauthorized, map[string]string{"erro": "Acesso expirado."})
				return
			}
			conta := contaAutenticada{ID: id, Nome: nome, Email: email, Perfil: perfil}
			if ultimo.Valid {
				conta.UltimoAcesso = &ultimo.Time
			}
			proximo.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), chaveConta{}, conta)))
		})
	}
}

func emitirSessao(w http.ResponseWriter, r *http.Request, b *sql.DB, segredo string, conta contaAutenticada, agora func() time.Time) error {
	segredoSessao := make([]byte, 32)
	if _, err := rand.Read(segredoSessao); err != nil {
		return err
	}
	renovacao := base64.RawURLEncoding.EncodeToString(segredoSessao)
	h := sha256.Sum256([]byte(renovacao))
	token, err := criarTokenAcesso(segredo, conta.ID, agora())
	if err != nil {
		return err
	}
	transacao, err := b.BeginTx(r.Context(), nil)
	if err != nil {
		return err
	}
	defer transacao.Rollback()
	if c, erroCookie := r.Cookie(nomeCookieSessao); erroCookie == nil {
		hashAnterior := sha256.Sum256([]byte(c.Value))
		resultado, erro := transacao.ExecContext(r.Context(), "UPDATE refresh_tokens SET revoked=1 WHERE token_hash=? AND revoked=0 AND expires_at>?", base64.RawURLEncoding.EncodeToString(hashAnterior[:]), agora().UTC())
		if erro != nil {
			return erro
		}
		linhas, erro := resultado.RowsAffected()
		if erro != nil {
			return erro
		}
		if linhas != 1 {
			return errors.New("sessão já renovada ou revogada")
		}
	}
	if _, err = transacao.ExecContext(r.Context(), "INSERT INTO refresh_tokens(user_id,token_hash,expires_at) VALUES(?,?,?)", conta.ID, base64.RawURLEncoding.EncodeToString(h[:]), agora().Add(duracaoTokenRenovacao).UTC()); err != nil {
		return err
	}
	if err = transacao.Commit(); err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{Name: nomeCookieSessao, Value: renovacao, Path: "/api/auth", HttpOnly: true, Secure: seguroRequisicao(r), SameSite: http.SameSiteStrictMode, Expires: agora().Add(duracaoTokenRenovacao), MaxAge: int(duracaoTokenRenovacao.Seconds())})
	w.Header().Set("Authorization", "Bearer "+token)
	w.Header().Set("Access-Control-Expose-Headers", "Authorization")
	return nil
}
func criarTokenAcesso(segredo string, id int64, agora time.Time) (string, error) {
	c := jwt.RegisteredClaims{Subject: strconv.FormatInt(id, 10), IssuedAt: jwt.NewNumericDate(agora), ExpiresAt: jwt.NewNumericDate(agora.Add(duracaoTokenAcesso)), Issuer: "mira"}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString([]byte(segredo))
}
func validarTokenAcesso(valor, segredo string, id int64, agora time.Time) bool {
	c := &jwt.RegisteredClaims{}
	t, e := jwt.ParseWithClaims(valor, c, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("algoritmo inválido")
		}
		return []byte(segredo), nil
	}, jwt.WithIssuer("mira"), jwt.WithTimeFunc(func() time.Time { return agora }))
	return e == nil && t.Valid && c.Subject == strconv.FormatInt(id, 10)
}
func seguroRequisicao(r *http.Request) bool {
	return r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")
}
func expirarCookieSessao(w http.ResponseWriter, seguro bool) {
	http.SetCookie(w, &http.Cookie{Name: nomeCookieSessao, Path: "/api/auth", HttpOnly: true, Secure: seguro, SameSite: http.SameSiteStrictMode, MaxAge: -1})
}
func responderJSON(w http.ResponseWriter, status int, valor any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(valor)
}
