# Arquitetura

## Protótipo atual

[`prototipo.html`](../prototipo.html) reúne a marcação das telas, o CSS dos temas e o JavaScript que controla a interface. O navegador executa tudo localmente, sem build, servidor ou bibliotecas externas.

| Parte | Responsabilidade |
| --- | --- |
| Marcação HTML | Define o login, a estrutura do console, as telas de projetos, containers e usuários e os diálogos. |
| CSS | Define os tokens claro e escuro, os componentes compartilhados, os estados visuais e o layout responsivo. |
| Estado em JavaScript | Mantém usuários, projetos, containers e histórico em memória. Mapas guardam os deploys ativos por slug e os logs por nome de container. |
| Navegação e renderização | `navigate()` alterna as telas no mesmo documento e coordena foco, filtros e ciclo dos logs. As funções de renderização atualizam as seções a partir do estado atual. |
| Simulação | Temporizadores geram linhas de log e avançam as etapas de deploy. Os dados da VPS e dos containers vêm da amostra local. |

Recarregar volta ao login e restaura os dados de exemplo. A única preferência gravada em `localStorage` é o tema, na chave `mira:tema`. A autenticação, as métricas e as operações são simuladas. A aplicação ainda não se comunica com a VPS nem com o Docker.

## MVP conectado previsto

```text
Navegador → https://mira.macb.tech → Nginx compartilhado
                                     └─ rede proxy → mira-http:6002
                                                       ├─ / → React compilado
                                                       └─ /api/ → Go
                                                                   ├─ SQLite em /srv/data/mira/
                                                                   ├─ rede privada → integração Docker
                                                                   ├─ métricas do host
                                                                   └─ GitHub Actions → deploy via Compose
```

- O serviço Go entrega o frontend React compilado com Vite e a API no mesmo domínio. O frontend recebe estatísticas, logs e progresso de deploy por Server-Sent Events (SSE).
- A API autentica o operador, aplica as regras de negócio e guarda usuários, sessões, projetos, deploys e auditoria no SQLite. O acesso às métricas da VPS usa `gopsutil`.
- O componente Docker fica em rede privada e concentra o acesso ao daemon. O inventário usa o label `mira.project=<slug>`. As operações de início, parada e reinício são liberadas por serviço após validação em containers de teste.
- O Nginx compartilhado termina o TLS com certificado gerenciado pelo Certbot. A porta `6002` é interna e corresponde à implantação inicial. A mudança posterior para `6000` segue a ordem documentada na especificação.
- O Mira solicita e acompanha o deploy no GitHub Actions. O workflow de cada projeto aplica a configuração do seu Compose na VPS.

A especificação detalha a [arquitetura de deploy](./especificacao.md#6-arquitetura-de-deploy), o [contrato da API](./especificacao.md#7-contrato-mínimo-da-api), os [dados persistidos](./especificacao.md#8-dados-persistidos) e a [estrutura prevista](./especificacao.md#9-estrutura-prevista) para `server/`, `web/` e `deploy/`.
