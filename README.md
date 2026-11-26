# Mira

Console web self-hosted para acompanhar e operar os containers Docker de uma VPS.

Abra [prototipo.html](./prototipo.html) no navegador e clique em **Entrar** com os campos já preenchidos. O protótipo usa dados simulados e funciona sem instalação.

## Desenvolvimento local

Use Go 1.26 e Node.js 24 LTS. O módulo Go seleciona a versão de correção usada no build. No PowerShell, use `npm.cmd` caso a política de execução impeça o uso de `npm`.

```sh
npm --prefix web ci
npm --prefix web run build
go -C server run ./cmd/api
```

Antes de iniciar a API pela primeira vez, defina no ambiente `MIRA_JWT_SECRET` com pelo menos 32 bytes aleatórios, `MIRA_ADMIN_EMAIL` e `MIRA_ADMIN_PASSWORD` com pelo menos 12 caracteres. Opcionalmente, defina `MIRA_ADMIN_NAME`. Em PowerShell, por exemplo:

```powershell
$env:MIRA_JWT_SECRET = [Convert]::ToHexString([Security.Cryptography.RandomNumberGenerator]::GetBytes(32))
$env:MIRA_ADMIN_NAME = 'Administrador'
$env:MIRA_ADMIN_EMAIL = 'admin@exemplo.com'
$env:MIRA_ADMIN_PASSWORD = 'defina-uma-senha-forte'
```

Use a senha real apenas no terminal local e não a salve no repositório. O administrador inicial é criado uma vez, apenas se o banco ainda não tiver usuários. Abra [Mira](http://127.0.0.1:6002) e consulte a [saúde da API](http://127.0.0.1:6002/api/saude). O Go serve a interface compilada e aplica as migrações automaticamente em `server/.local/mira.db`.

Para editar a interface com recarga automática, mantenha o Go em execução e rode `npm --prefix web run dev` em outro terminal. O Vite encaminha `/api` ao Go local. O Docker é necessário para testar o inventário com o Compose. Veja a [preparação para a VPS](./deploy/README.md).

No Windows, quando o Controle de Aplicativo bloquear o executável gerado pelo Go, rode a API e a interface compilada no Docker Desktop com containers Linux. A porta abaixo fica restrita ao localhost e o SQLite fica no volume `mira-dev-data`. O mesmo servidor Vite pode encaminhar as chamadas a esse ambiente.

```sh
docker build --target http -t mira-dev .
docker run --rm --name mira-dev -p 127.0.0.1:6002:6002 -v mira-dev-data:/srv/data/mira mira-dev
```

[Documentação](./docs/README.md) · [Plano de implementação](./PLANO.md)
