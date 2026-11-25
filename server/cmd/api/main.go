package main

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dev-macb/mira/server/api"
	"github.com/dev-macb/mira/server/internal/docker"
	servidorhttp "github.com/dev-macb/mira/server/internal/httpserver"
	armazenamento "github.com/dev-macb/mira/server/internal/store"
)

func main() {
	if erro := executarAPI(); erro != nil {
		log.Fatal(erro)
	}
}

func executarAPI() error {
	contextoProcesso, pararMonitoramentoSinais := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer pararMonitoramentoSinais()

	var clienteInventario *docker.ClienteInventario
	if urlComponenteDocker := os.Getenv("MIRA_DOCKER_URL"); urlComponenteDocker != "" {
		var erro error
		clienteInventario, erro = docker.NovoClienteInventario(urlComponenteDocker, os.Getenv("MIRA_DOCKER_TOKEN"))
		if erro != nil {
			return erro
		}
	}

	contextoInicializacao, cancelarInicializacao := context.WithTimeout(contextoProcesso, 30*time.Second)
	caminhoBanco := cmp.Or(os.Getenv("MIRA_DB_PATH"), ".local/mira.db")
	bancoDeDados, erro := armazenamento.AbrirSQLite(contextoInicializacao, caminhoBanco)
	cancelarInicializacao()
	if erro != nil {
		return fmt.Errorf("abrir SQLite: %w", erro)
	}
	defer bancoDeDados.Close()

	enderecoHTTP := cmp.Or(os.Getenv("MIRA_HTTP_ADDR"), "127.0.0.1:6002")
	diretorioInterface := cmp.Or(os.Getenv("MIRA_WEB_DIR"), "../web/dist")
	roteador := api.NovoRoteador(bancoDeDados, clienteInventario, diretorioInterface)

	log.Printf("Mira em %s", enderecoHTTP)
	return servidorhttp.Servir(contextoProcesso, enderecoHTTP, roteador)
}
