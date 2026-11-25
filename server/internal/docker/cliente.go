package docker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const rotaInventario = "/inventario"

type ContainerVinculado struct {
	SlugProjeto string `json:"project"`
	Nome        string `json:"name"`
	Imagem      string `json:"image"`
	Estado      string `json:"state"`
}

type ClienteInventario struct {
	urlInventario string
	tokenAcesso   string
	clienteHTTP   *http.Client
}

func NovoClienteInventario(urlBase, tokenAcesso string) (*ClienteInventario, error) {
	if erro := ValidarTokenAcesso(tokenAcesso); erro != nil {
		return nil, erro
	}
	return &ClienteInventario{
		urlInventario: strings.TrimRight(urlBase, "/") + rotaInventario,
		tokenAcesso:   tokenAcesso,
		clienteHTTP: &http.Client{
			Timeout:   3 * time.Second,
			Transport: &http.Transport{IdleConnTimeout: 30 * time.Second},
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}, nil
}

func (inventario *ClienteInventario) ListarContainers(contextoConsulta context.Context) ([]ContainerVinculado, error) {
	requisicao, erro := http.NewRequestWithContext(contextoConsulta, http.MethodGet, inventario.urlInventario, nil)
	if erro != nil {
		return nil, erro
	}
	requisicao.Header.Set("Authorization", "Bearer "+inventario.tokenAcesso)
	resposta, erro := inventario.clienteHTTP.Do(requisicao)
	if erro != nil {
		return nil, erro
	}
	defer resposta.Body.Close()
	if resposta.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("componente Docker retornou HTTP %d", resposta.StatusCode)
	}
	var containersVinculados []ContainerVinculado
	if erro := json.NewDecoder(resposta.Body).Decode(&containersVinculados); erro != nil {
		return nil, erro
	}
	return containersVinculados, nil
}

func ValidarTokenAcesso(tokenAcesso string) error {
	if len(strings.TrimSpace(tokenAcesso)) < 32 || strings.ContainsAny(tokenAcesso, "\r\n\t ") {
		return errors.New("MIRA_DOCKER_TOKEN deve ter ao menos 32 caracteres e ser informado sem espaços")
	}
	return nil
}
