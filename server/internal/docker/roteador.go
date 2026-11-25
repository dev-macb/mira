package docker

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/moby/moby/client"
)

const rotuloProjeto = "mira.project"

func NovoRoteadorInventario(clienteDocker *client.Client, tokenAcesso string) http.Handler {
	autorizacaoEsperada := []byte("Bearer " + tokenAcesso)
	roteador := chi.NewRouter()
	roteador.Get(rotaInventario, func(gravadorResposta http.ResponseWriter, requisicao *http.Request) {
		autorizacao := requisicao.Header.Get("Authorization")
		if subtle.ConstantTimeCompare([]byte(autorizacao), autorizacaoEsperada) != 1 {
			http.Error(gravadorResposta, "Acesso não autorizado.", http.StatusUnauthorized)
			return
		}

		contextoListagem, cancelarListagem := context.WithTimeout(requisicao.Context(), 2*time.Second)
		defer cancelarListagem()

		resultadoListagem, erro := clienteDocker.ContainerList(contextoListagem, client.ContainerListOptions{
			All:     true,
			Filters: make(client.Filters).Add("label", rotuloProjeto),
		})
		if erro != nil {
			log.Printf("listar containers: %v", erro)
			http.Error(gravadorResposta, "Inventário indisponível.", http.StatusBadGateway)
			return
		}

		containersVinculados := make([]ContainerVinculado, 0, len(resultadoListagem.Items))
		for _, containerDocker := range resultadoListagem.Items {
			slugProjeto := containerDocker.Labels[rotuloProjeto]
			if slugProjeto == "" || len(containerDocker.Names) == 0 {
				continue
			}
			containersVinculados = append(containersVinculados, ContainerVinculado{
				SlugProjeto: slugProjeto,
				Nome:        strings.TrimPrefix(containerDocker.Names[0], "/"),
				Imagem:      containerDocker.Image,
				Estado:      string(containerDocker.State),
			})
		}

		gravadorResposta.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(gravadorResposta).Encode(containersVinculados)
	})
	return roteador
}
