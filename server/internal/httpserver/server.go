package httpserver

import (
	"context"
	"errors"
	"net/http"
	"time"
)

func Servir(contextoProcesso context.Context, enderecoHTTP string, tratadorHTTP http.Handler) error {
	servidorHTTP := &http.Server{
		Addr:              enderecoHTTP,
		Handler:           tratadorHTTP,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	erroAoServir := make(chan error, 1)
	go func() {
		erroAoServir <- servidorHTTP.ListenAndServe()
	}()

	select {
	case erroServidor := <-erroAoServir:
		return erroServidor
	case <-contextoProcesso.Done():
	}

	contextoEncerramento, cancelarEncerramento := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelarEncerramento()

	if erro := servidorHTTP.Shutdown(contextoEncerramento); erro != nil {
		return erro
	}
	if erroServidor := <-erroAoServir; !errors.Is(erroServidor, http.ErrServerClosed) {
		return erroServidor
	}
	return nil
}
