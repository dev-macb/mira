package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dev-macb/mira/server/internal/docker"
)

func tratarSaude(bancoDeDados *sql.DB, clienteInventario *docker.ClienteInventario) http.HandlerFunc {
	return func(gravadorResposta http.ResponseWriter, requisicao *http.Request) {
		contextoVerificacao, cancelarVerificacao := context.WithTimeout(requisicao.Context(), 3*time.Second)
		defer cancelarVerificacao()

		codigoHTTP, estadoSaude := http.StatusOK, "ok"
		if erro := verificarDependencias(contextoVerificacao, bancoDeDados, clienteInventario); erro != nil {
			log.Print(erro)
			codigoHTTP, estadoSaude = http.StatusServiceUnavailable, "indisponivel"
		}

		gravadorResposta.Header().Set("Content-Type", "application/json; charset=utf-8")
		gravadorResposta.Header().Set("Cache-Control", "no-store")
		gravadorResposta.WriteHeader(codigoHTTP)
		_ = json.NewEncoder(gravadorResposta).Encode(map[string]string{"status": estadoSaude})
	}
}

func verificarDependencias(contextoVerificacao context.Context, bancoDeDados *sql.DB, clienteInventario *docker.ClienteInventario) error {
	if erro := bancoDeDados.PingContext(contextoVerificacao); erro != nil {
		return fmt.Errorf("saúde SQLite: %w", erro)
	}
	if clienteInventario != nil {
		if _, erro := clienteInventario.ListarContainers(contextoVerificacao); erro != nil {
			return fmt.Errorf("saúde do inventário Docker: %w", erro)
		}
	}
	return nil
}
