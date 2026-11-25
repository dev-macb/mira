package api

import (
	"database/sql"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dev-macb/mira/server/internal/docker"
	"github.com/go-chi/chi/v5"
)

func NovoRoteador(bancoDeDados *sql.DB, clienteInventario *docker.ClienteInventario, diretorioInterface string) http.Handler {
	roteador := chi.NewRouter()
	roteador.Get("/api/saude", tratarSaude(bancoDeDados, clienteInventario))

	caminhoPaginaInicial := filepath.Join(diretorioInterface, "index.html")
	roteador.Get("/", func(gravadorResposta http.ResponseWriter, requisicao *http.Request) {
		if _, erro := os.Stat(caminhoPaginaInicial); os.IsNotExist(erro) {
			http.Error(gravadorResposta, "Compile a interface com npm run build ou abra o Vite.", http.StatusServiceUnavailable)
			return
		}

		http.ServeFile(gravadorResposta, requisicao, caminhoPaginaInicial)
	})

	servidorArquivosEstaticos := http.FileServer(http.Dir(filepath.Join(diretorioInterface, "assets")))
	roteador.Handle("/assets/*", http.StripPrefix("/assets/", servidorArquivosEstaticos))

	return roteador
}
