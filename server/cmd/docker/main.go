package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dev-macb/mira/server/internal/docker"
	servidorhttp "github.com/dev-macb/mira/server/internal/httpserver"
	"github.com/moby/moby/client"
)

func main() {
	if erro := executarInventario(); erro != nil {
		log.Fatal(erro)
	}
}

func executarInventario() error {
	tokenAcesso := os.Getenv("MIRA_DOCKER_TOKEN")
	if erro := docker.ValidarTokenAcesso(tokenAcesso); erro != nil {
		return erro
	}
	clienteDocker, erro := client.New()
	if erro != nil {
		return erro
	}
	defer clienteDocker.Close()

	contextoProcesso, pararMonitoramentoSinais := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer pararMonitoramentoSinais()

	roteador := docker.NovoRoteadorInventario(clienteDocker, tokenAcesso)
	return servidorhttp.Servir(contextoProcesso, "0.0.0.0:8080", roteador)
}
