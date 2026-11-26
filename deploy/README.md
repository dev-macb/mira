# Implantação inicial na VPS

O Mira usa a rede externa `proxy`, o upstream `mira-http:6002` e os dados em `/srv/data/mira`. O Nginx e o Certbot continuam em `/srv/infra`. O Compose do Mira cria a rede privada do inventário e inicia os dois serviços sem publicar portas.

## Preparação

Confirme que o registro A de `mira.macb.tech` aponta para o IPv4 público da VPS, `187.77.63.198`. Confira se `6002` está disponível no inventário interno, se a rede externa `proxy` está ativa e se `/srv/data/certbot/www` e `/srv/data/certbot/conf` estão montados pelo Certbot compartilhado. A porta do Mira é interna e não precisa estar publicada no host.

```sh
sudo install -d -o 10001 -g 1001 -m 0750 /srv/data/mira
cd /srv/apps/mira
cp .env.example .env
# Gere e configure os segredos e a conta inicial no arquivo .env com acesso restrito.
stat -c '%g' /var/run/docker.sock
openssl rand -hex 32
openssl rand -hex 32
docker compose config --quiet
docker compose up -d --build --wait
```

O serviço HTTP roda com UID `10001` e recebe acesso de escrita ao diretório de dados pelo grupo `1001` do host. Somente `docker-inventory` monta o socket e recebe seu GID real em `MIRA_DOCKER_GID`. Gere valores aleatórios separados para `MIRA_DOCKER_TOKEN` e `MIRA_JWT_SECRET`; configure também `MIRA_ADMIN_EMAIL` e uma senha inicial de pelo menos 12 caracteres em `MIRA_ADMIN_PASSWORD`. O administrador inicial só é criado quando o banco ainda não tem usuários. Para a aplicação local, faça uma cópia autorizada do checkout para `/srv/apps/mira`. Siga o procedimento interno de publicação para transferir os arquivos. A conta SSH deste ambiente não permite executar o fluxo completo de `sudo`, configurar o dono do diretório e publicar automaticamente cada atualização na VPS.

## Nginx e certificado

Ative primeiro [mira.http.conf](./nginx/mira.http.conf) em `/srv/infra/nginx/conf.d/mira.conf`, valide com `cd /srv/infra && docker compose exec -T nginx nginx -t` e aplique com `docker compose exec -T nginx nginx -s reload`. Para verificar o acesso ao desafio antes de emitir o certificado, crie um arquivo de prova no webroot e solicite-o por HTTP com o host `mira.macb.tech`.

Emita o certificado pelo Certbot compartilhado em modo webroot:

```sh
cd /srv/infra
docker compose --profile maintenance run --rm --no-deps certbot certonly \
  --config /etc/certbot/cli.ini \
  --webroot -w /var/www/certbot \
  --cert-name mira.macb.tech \
  -d mira.macb.tech
```

Depois substitua `/srv/infra/nginx/conf.d/mira.conf` pela versão TLS de [mira.https.conf](./nginx/mira.https.conf), valide novamente com `nginx -t`, aplique a recarga graciosa e confirme `https://mira.macb.tech`. O upstream `mira-http:6002` integra-se à resolução dinâmica já habilitada no Nginx compartilhado.

**Atenção à renovação:** no servidor deste deploy os arquivos `infra-certbot-renew.service` e `infra-certbot-renew.timer` estão em `/srv/infra/certbot`, mas o timer não está instalado/ativo no systemd. A renovação e a recarga do Nginx precisam ser habilitadas por um administrador do servidor antes do vencimento do certificado. Depois, confirme com `systemctl list-timers --all`.

Identifique o IP ou CIDR autorizado do Nginx para configurar `TRUST_PROXY` na etapa de autenticação. A base da API ainda não usa cabeçalhos encaminhados para identificar clientes.

## Validação

- Confirme `healthy` em `docker compose ps` e confira a coluna de portas. A verificação HTTP exercita o SQLite e a chamada autenticada ao inventário pela rede privada.
- Com containers de teste, confira se a listagem privada inclui um container com `mira.project=<slug>`, inclusive quando parado, e se omite containers sem o label. A partir do checkout, a consulta de diagnóstico abaixo usa o token da própria API, sem imprimi-lo.

  ```sh
  docker compose exec -T http sh -c 'wget -qO- --header="Authorization: Bearer $MIRA_DOCKER_TOKEN" "$MIRA_DOCKER_URL/inventario"'
  ```

- Confirme o redirecionamento de HTTP para HTTPS e teste a página e `https://mira.macb.tech/api/saude` com o certificado público. Recrie o serviço HTTP e confirme que o SQLite continua disponível.
