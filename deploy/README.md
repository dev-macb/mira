# Implantação inicial na VPS

O Mira usa a rede externa `proxy`, o upstream `mira-http:6002` e os dados em `/srv/data/mira`. O Nginx e o Certbot continuam em `/srv/infra`. O Compose do Mira cria a rede privada do inventário e inicia os dois serviços sem publicar portas.

## Preparação

Confira se o DNS de `mira.macb.tech` aponta para a VPS, se `6002` está disponível no inventário interno e se a rede `proxy` corresponde à do Nginx compartilhado. No checkout em `/srv/apps/mira`, crie um `.env` a partir de [`.env.example`](../.env.example), caso ainda não exista. Preencha `MIRA_DOCKER_GID` com `stat -c %g /var/run/docker.sock` e `MIRA_DOCKER_TOKEN` com a saída de `openssl rand -hex 32`.

```sh
docker network inspect proxy
sudo install -d -o 10001 -g 10001 -m 0750 /srv/data/mira
docker compose config --quiet
docker compose up -d --build --wait
```

O serviço HTTP roda como UID/GID `10001:10001`. Somente `docker-inventory` monta o socket e recebe o GID do host. O modo de consulta é definido pelo serviço, que só chama a listagem do Docker e exige o token da API. O acesso ao socket continua privilegiado no daemon.

## Nginx e certificado

Use o fluxo de ativação já adotado em `/srv/infra`. Confira se o webroot compartilhado aparece como `/var/www/certbot` no Nginx. Ative primeiro [mira.http.conf](./nginx/mira.http.conf) para servir o desafio ACME, valide com `nginx -t`, recarregue e emita o certificado pelo Certbot compartilhado com `--webroot -w /var/www/certbot -d mira.macb.tech`. Depois substitua o vhost ativo por [mira.https.conf](./nginx/mira.https.conf) e repita a validação e a recarga. O vhost resolve o alias pela rede Docker.

Identifique o IP ou CIDR autorizado do Nginx para configurar `TRUST_PROXY` na etapa de autenticação. A base da API ainda não usa cabeçalhos encaminhados para identificar clientes.

## Validação

- Confirme `healthy` em `docker compose ps` e confira a coluna de portas. A verificação HTTP exercita o SQLite e a chamada autenticada ao inventário pela rede privada.
- Com containers de teste, confira se a listagem privada inclui um container com `mira.project=<slug>`, inclusive quando parado, e se omite containers sem o label. A partir do checkout, a consulta de diagnóstico abaixo usa o token da própria API, sem imprimi-lo.

  ```sh
  docker compose exec -T http sh -c 'wget -qO- --header="Authorization: Bearer $MIRA_DOCKER_TOKEN" "$MIRA_DOCKER_URL/inventario"'
  ```

- Confirme o redirecionamento de HTTP para HTTPS e teste a página e `https://mira.macb.tech/api/saude` com o certificado público. Recrie o serviço HTTP e confirme que o SQLite continua disponível.
