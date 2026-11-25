# Mira — requisitos do MVP

O Mira é um console web self-hosted para administradores acompanharem e operarem os containers Docker de uma VPS. **Projetos** reúne as métricas da VPS e o inventário de projetos. O detalhe do projeto reúne recursos, containers e deploys. Abrir um container dá acesso às suas métricas, ações e logs. O menu da conta dá acesso à gestão de usuários, ao tema e à saída.

[**Abrir o protótipo Mira**](../prototipo.html). O HTML funciona localmente com dados simulados em memória. Recarregar restaura a amostra e retorna ao login. A preferência de tema persiste no navegador. O protótipo serve de referência para a interface. Esta especificação descreve a autenticação, a persistência, as operações na VPS e o deploy via GitHub Actions do MVP conectado.

## 1. Objetivo e escopo

Entregar acesso por HTTPS em um subdomínio dedicado, autenticação e gestão de contas administrativas, cadastro de projetos, métricas da VPS e dos containers, logs, início, parada e reinício de containers, deploy manual com acompanhamento e histórico, e auditoria. O vínculo entre projeto e container usa o label Docker `mira.project=<slug>`.

O backlog reúne 2FA, allowlist de IP/VPN, perfis com permissões distintas, isolamento por projeto, execução de comandos em containers, editor de `.env`, Git, webhooks de CI/CD, backup de volumes, console PostgreSQL, alertas, healthchecks personalizados, gráficos históricos, rollback automático, agendamento e botão de pânico.

## 2. Stack

| Camada | Tecnologia e responsabilidade |
| --- | --- |
| Backend | Go com chi. Entrega a interface e a API, autentica, aplica as regras de operação, consulta métricas e publica eventos SSE. |
| Frontend | React, Vite e CSS próprio. O serviço Go entrega os arquivos compilados no mesmo domínio da API. |
| Banco | SQLite em `/srv/data/mira/mira.db`. Guarda usuários, tokens de renovação, projetos, deploys e auditoria. |
| Proxy reverso | Nginx compartilhado em `/srv/infra`, com TLS gerenciado pelo Certbot. Encaminha o tráfego ao Mira pela rede `proxy`. |
| Acesso ao Docker | Componente dedicado em rede privada, com Docker SDK e acesso restrito às operações do MVP. |
| Deploy | GitHub Actions aciona o procedimento de deploy via SSH e Docker Compose de cada projeto autorizado. |
| Métricas | `gopsutil` com acesso às métricas do host. CPU, RAM, disco, load e uptime representam a VPS. |

O protótipo local reúne HTML, CSS e JavaScript e usa fontes do sistema.

## 3. Telas e comportamentos

### 3.1 Acesso ao console

- O painel translúcido fica centralizado e tem largura máxima de 370 px. A mira laranja em caracteres Braille Unicode aparece acima dos campos **E-mail** e **Senha**, com rótulos sempre visíveis. O botão **Entrar** ocupa toda a largura do formulário.
- O fundo fixo de todas as telas combina a grade de bússolas com a textura pontilhada. Cada bússola mede 38 × 3 px, com intervalo de 56 px entre centros. Abaixo de 650 px, mede 30 × 3 px, com intervalo de 44 px. A ponta azul segue o cursor ou o toque e a outra ponta é laranja nos dois temas. Antes da primeira interação, as barras apontam para o centro. Com movimento reduzido, o fundo permanece estático durante a interação.
- **Mostrar/Ocultar** alterna a visibilidade da senha. Enter envia o formulário. O login usa apenas contas ativas. Credenciais válidas atualizam **Último acesso**, mostram o nome da conta no cabeçalho e abrem o topo de **Projetos**.
- Uma falha mantém os valores, sinaliza os campos, exibe **E-mail ou senha incorretos** e devolve o foco ao e-mail. O alerta some ao editar. No HTML local, os campos vêm preenchidos para demonstração e **Restaurar campos** recupera os valores iniciais, oculta a senha e foca o e-mail.
- No MVP conectado, o formulário começa vazio, indica o envio em andamento e bloqueia envios duplicados. **Sair da conta** encerra a sessão, fecha o menu e retorna ao login com foco no e-mail. No protótipo, sair reinicia a navegação e preserva as alterações em memória até a recarga.

### 3.2 Projetos

- Abaixo de **Projetos** e **Novo projeto**, o cabeçalho exibe a VPS no formato `domínio - IP` e o tempo de atividade. Os quatro cartões aparecem na ordem **Containers ativos**, **CPU**, **Memória** e **Disco**. Containers ativos mostra o total como denominador. CPU mostra percentual, vCPUs e load. Memória e disco mostram percentual, uso e capacidade em GB. As barras de CPU e memória são azuis. A de disco é laranja.
- A tabela **Projetos** fica abaixo das métricas da VPS. As colunas são **Projeto**, **Containers ativos**, **Situação**, **Cliente** e **Ações**. O domínio fica abaixo do nome. O indicador de situação é azul para **Em produção**, laranja para **Em desenvolvimento** e neutro para **Arquivado**. O lápis abre a edição e a linha ou o botão com o nome abre o detalhe.
- **Buscar projeto** encontra trechos do nome, slug, cliente, domínio ou nome de qualquer container vinculado, sem distinção entre maiúsculas e minúsculas. O contador acompanha a busca. A tabela distingue busca sem resultados de ausência de projetos. **Novo projeto** fica disponível na lista vazia.

### 3.3 Cadastro de projetos

- **Novo projeto** abre um diálogo com **Nome do projeto** (até 80 caracteres), **Slug** (63), **Cliente** (80), **Situação**, **Domínio** (253) e **Descrição** (2.000). Nome e slug são obrigatórios e o slug é único. O formato aceita letras minúsculas e números, com hífens simples entre segmentos. O texto de ajuda atualiza o label `mira.project=<slug>` enquanto o operador digita e usa `meu-projeto` como exemplo inicial.
- Cliente vazio assume **Pessoal**. Domínio ausente aparece como **Sem domínio** na tabela. A situação inicial é **Em desenvolvimento**. O cadastro também oferece **Em produção** e **Arquivado**.
- **Criar projeto** abre o detalhe do novo projeto e exibe **Projeto criado.**. Um nome composto apenas de espaços ou um slug duplicado mantém o diálogo aberto e sinaliza o campo correspondente.
- O lápis na lista abre **Editar projeto** com os valores atuais e mantém o slug em modo de leitura. **Salvar alterações** atualiza a linha e o detalhe, preserva a busca e devolve o foco ao lápis da linha ou à busca se a linha sair do resultado. Em um projeto já cadastrado, o arquivamento passa pela confirmação no detalhe. A edição permite reativar um projeto arquivado.

### 3.4 Detalhe do projeto

- O caminho **MIRA / PROJETOS / nome do projeto** fica acima do título. O cabeçalho mostra o nome do projeto, a situação imediatamente à direita e as ações **Arquivar projeto** e **Fazer deploy**. Abaixo aparecem o cliente, o domínio quando cadastrado e a descrição quando preenchida, com as quebras de linha preservadas.
- **Recursos do projeto** ocupa toda a largura antes do inventário. Exibe, nesta ordem, ativos sobre vinculados, soma de CPU dos ativos, memória em uso dos ativos sobre a soma dos limites de todos os vinculados e maior uptime dos ativos com o nome do container. Com containers vinculados e todos inativos, CPU e memória em uso mostram zero e o uptime mostra `—`. Sem vínculos, a contagem é zero e as demais leituras mostram `—`.
- A tabela **Containers** tem as colunas **Container**, **Estado**, **Memória / limite** e **Ações**. **Buscar container** combina busca por trecho do nome com o seletor **Todos os estados**, `running`, `stopped`, `exited` ou `unhealthy`. A linha ou o botão com o nome abre o detalhe. Projeto sem containers oculta os filtros, orienta a usar `mira.project=<slug>` e desabilita **Fazer deploy**. Filtros sem resultados têm mensagem própria. No MVP, **Fazer deploy** também depende de um workflow habilitado para o projeto.
- **Histórico de deploys** vem depois dos containers e mostra **Versão**, **Quando** e **Resultado**. No protótipo, a versão é a tag e o alvo é o container. No MVP, a versão é o SHA abreviado e o histórico identifica os serviços do Compose afetados e fornece o link da execução. O seletor oferece **Todos os resultados**, **Concluídos** e **Falhos**. Os resultados aparecem como **Concluído** em azul ou **Falhou** em vermelho. O resultado novo entra no topo. Histórico vazio e filtro sem correspondência têm mensagens distintas.
- Voltar do container ao mesmo projeto preserva filtros e páginas. Mudar de projeto limpa os filtros de containers e de deploys e volta à primeira página dessas tabelas. As leituras refletem as ações de container e o resultado dos deploys.
- **Arquivar projeto** pede confirmação e muda somente a situação. O projeto segue consultável e conserva containers e histórico, com a ação de arquivar desabilitada. As ações de container e o deploy continuam disponíveis conforme os vínculos e estados do projeto.

### 3.5 Detalhe do container e logs

- O cabeçalho mostra o nome, o estado Docker ao lado e a imagem no formato `imagem:tag` abaixo. Os cartões exibem **CPU em uso**, **Memória em uso** com limite em MB e **Uptime**. Abaixo de uma hora, o uptime aparece como `< 1 h`. Para inativos, o limite continua visível e as leituras indisponíveis exibem `—`. O caminho **MIRA / PROJETOS / nome do projeto / NOME DO CONTAINER** oferece retorno à lista e ao projeto.
- `running` e `unhealthy` contam como ativos. A interface mostra `running` em azul, `unhealthy` em laranja e `stopped` e `exited` em vermelho. **Reiniciar** e **Parar** aparecem para ativos, nessa ordem, e **Iniciar** aparece para inativos. No MVP, as ações reais dependem de autorização explícita por serviço do Compose. A tabela usa ícones e o detalhe usa ícones com texto.
- A confirmação identifica ação, alvo e impacto. Cancelar mantém o estado. Uma operação bem-sucedida atualiza o estado, o uptime, as métricas e os totais afetados. O protótipo também acrescenta uma linha `INFO` ao buffer e mostra uma notificação com o nome do container e a ação concluída.
- **Logs do container** é a última seção e abre automaticamente com o filtro limpo e a tela no fim do buffer. Sair encerra o fluxo. Cada container mantém as últimas 200 linhas, descartando as mais antigas. Cada linha mostra hora no formato local de 24 horas, nível (`INFO`, `WARN` ou `ERR`) e mensagem.
- **Buscar nos logs** pesquisa trechos da mensagem e do nível, sem distinguir maiúsculas de minúsculas, nas linhas armazenadas e nas novas. Filtrar preserva o buffer e os horários. O rodapé mostra o total de linhas ou a quantidade filtrada em relação ao total. Uma mensagem orienta a ajustar a busca quando nada corresponde.
- O estado indica **Recebendo logs**, **Pausado** ou **Container inativo**. **Pausar** congela a leitura e **Retomar** volta ao fluxo. Um container inativo mantém as últimas linhas disponíveis para consulta e desabilita o controle de pausa. Novas linhas levam a rolagem ao fim apenas quando o leitor já está perto do fim. Alterar a busca ou alternar a pausa leva a visualização ao fim das linhas exibidas.
- No MVP conectado, abrir o detalhe inicia o carregamento dos logs. Pausar ou sair fecha a conexão SSE. Retomar ou voltar a um container ativo abre a conexão correspondente. A interface informa carregamento, desconexão e reconexão, preservando as linhas recebidas.

### 3.6 Deploy e confirmações

- No protótipo, **Fazer deploy** abre o formulário **Container de destino**, **Imagem** e **Tag**. Trocar o destino carrega a imagem e a tag atuais. Ao iniciar, a simulação exibe o alvo, `imagem:tag`, a situação e as etapas de download, parada, recriação, início e verificação de saúde. A tag `indisponivel` simula falha antes de alterar o container anterior.
- No MVP conectado, **Fazer deploy** mostra a revisão aprovada, os serviços abrangidos e o workflow associado ao projeto. Ao confirmar, o Mira solicita a execução no GitHub Actions e exibe seu andamento e o link para os logs da execução. O workflow aplica o Compose e verifica a saúde dos serviços.
- Cada projeto aceita um deploy em andamento. **Acompanhar deploy** reabre o progresso. Fechar o diálogo ou mudar de tela mantém a operação. Início, parada e reinício ficam bloqueados para os alvos afetados durante a operação.
- A conclusão registra um único resultado no histórico do projeto de origem e aparece no diálogo aberto ou em uma notificação. No MVP, a interface apresenta a etapa que falhou quando disponível e atualiza o inventário com o estado observado dos containers. Ao concluir, a ação principal volta a **Fazer deploy**.
- Diálogos aceitam Escape e o botão **Cancelar** ou **Fechar**, conforme a etapa. As confirmações de arquivamento, gestão de usuários e ações de container focam **Cancelar** ao abrir.

### 3.7 Usuários

- **Usuários** no menu abre **MIRA / USUÁRIOS**. O painel **Contas do console** mostra **Nome**, **E-mail**, **Perfil**, **Situação**, **Último acesso** e **Ações**. O perfil é **Administrador**. A situação aparece como **Ativo** em azul ou **Inativo** em tom neutro. O último acesso vem do login bem-sucedido e aparece em formato `pt-BR` ou como `—` para conta sem acesso registrado.
- **Buscar usuário** pesquisa trechos do nome, e-mail ou perfil sem distinção entre maiúsculas e minúsculas. A tabela exibe **Nenhum usuário corresponde à busca.** quando nenhum registro corresponde.
- **Novo usuário** abre um diálogo com nome obrigatório de até 100 caracteres, e-mail válido de até 254 caracteres e senha obrigatória. O e-mail é normalizado para minúsculas e deve ser único. A nova conta é criada como ativa. Validações de nome ou e-mail mantêm o diálogo e indicam o campo.
- Na linha, o lápis abre **Editar usuário** com nome e e-mail atuais. **Nova senha (opcional)** começa vazia. Salvar com esse campo vazio conserva a senha. Editar o próprio nome atualiza o cabeçalho imediatamente. Ao fechar, o diálogo limpa o campo de senha.
- O ícone de energia aciona **Ativar** ou **Desativar** e pede confirmação com o e-mail e o efeito sobre o acesso. O controle fica desabilitado na conta em uso, com a indicação **Conta em uso**. Criar, editar, ativar e desativar atualizam a tabela, exibem uma notificação e devolvem o foco ao lápis da conta ou à busca quando a conta fica fora da página visível. No protótipo, contas e último acesso duram até a recarga.

## 4. Requisitos funcionais da implementação

| ID | Requisito | Resultado observável |
| --- | --- | --- |
| RF-01 | Autenticar contas administrativas ativas com bcrypt, rate limiting, JWT de 15 min e refresh rotativo persistido. | Login bem-sucedido cria sessão e atualiza o último acesso. Renovação substitui o refresh anterior. Logout revoga a sessão. |
| RF-02 | Criar, consultar, editar, arquivar e reativar projetos. Validar campos, unicidade e imutabilidade do slug na API. | Os dados sobrevivem à recarga. Conflito de slug aparece no campo. Arquivar preserva vínculos, estado dos containers e histórico. |
| RF-03 | Resolver vínculos por `mira.project` e fornecer nome, imagem, tag, estado, CPU, memória, limite e uptime. | Os totais e agregados correspondem aos containers do projeto e tratam `unhealthy` como ativo. |
| RF-04 | Executar `start`, `stop` e `restart` em serviços autorizados, validando estado, vínculo e concorrência com deploy. | O operador recebe o estado confirmado ou uma falha com nome do alvo e ação. O resultado é auditado. |
| RF-05 | Ler domínio e IP da VPS, vCPUs, CPU, memória, disco, load e uptime. | O painel indica unidades e capacidade, mantém a última leitura durante uma reconexão e a identifica como desatualizada. |
| RF-06 | Fornecer o trecho recente e o fluxo de logs do container aberto por SSE a partir de `ContainerLogs`, com `follow=true` para fonte ativa. | O fluxo respeita pausa e navegação. Busca, níveis e horários funcionam em um buffer de até 200 linhas por container. |
| RF-07 | Acionar manualmente o workflow autorizado para uma revisão aprovada. O workflow executa o deploy via Compose e verifica a saúde. | A API valida projeto e revisão, aceita um deploy em andamento por projeto e impede ações conflitantes nos serviços afetados. Uma falha anterior à aplicação do Compose preserva os containers em execução. |
| RF-08 | Consultar o histórico e recuperar o progresso por projeto. Persistir as etapas e o resultado final. | Sair do diálogo ou mudar de tela mantém a operação. O acompanhamento recupera as etapas já registradas. Cada operação gera um único registro com resultado. |
| RF-09 | Listar, criar, editar, ativar e desativar contas administrativas. Validar unicidade do e-mail normalizado e bloquear a desativação da conta autenticada na API. | A tabela exibe situação e último acesso. A senha é exigida ao criar e opcional ao editar. Conta inativa perde acesso. As alterações são auditadas. |

## 5. Identidade e componentes

**Mira** é o nome de apresentação. A aba usa **Mira · Acesso ao console** no login e **nome da tela · Mira** após entrar. `mira` é o prefixo técnico da rede, dos serviços e do label. O favicon e o atalho do cabeçalho usam a mesma mira pixelada em uma grade de 16 × 16, com contorno e quatro braços azuis e centro laranja. O atalho mede 32 × 32 px. O login usa a arte em caracteres Braille e disponibiliza o nome Mira para leitores de tela.

O layout segue o Breve Link. A página usa Courier New, fundo pontilhado, coluna central de até 960 px, cabeçalho horizontal e rodapé. O recuo lateral é de 48 px, passa a 32 px até 900 px e a 24 px até 650 px. O rodapé mostra `Mira // domínio - IP` e `[ @dev-macb ]`. Os títulos de página terminam em um ponto laranja e os títulos de painel e de seção usam colchetes. Os controles têm cantos retos. As métricas da VPS e do projeto passam de quatro para duas colunas até 900 px. As do container passam de três colunas para linhas empilhadas até 550 px.

| Cor de referência | Escuro | Claro |
| --- | --- | --- |
| Fundo | `#090f1e` | `#f1f6ff` |
| Superfície | `#101b2e` | `#ffffff` |
| Texto | `#ecf4ff` | `#101c31` |
| Destaque azul | `#38b6ff` | `#164bde` |
| Botão principal | `#38b6ff` com texto `#041427` | `#164eff` com texto `#ffffff` |
| Laranja | `#ff791f` | `#f55000` |
| Aviso | `#ff944d` | `#963600` |
| Erro e ação destrutiva | `#ff8097` | `#ad193c` |

O cabeçalho mostra a mira à esquerda e o nome e perfil da conta à direita junto ao menu. Até 650 px, ficam visíveis a mira e o menu. A mira abre o topo de **Projetos** com foco no título da página. No menu, **Projetos** vai direto à tabela e foca seu título. O menu também traz **Usuários**, o interruptor **Tema escuro** e **Sair da conta**. Clique fora ou foco fora fecha o menu. Escape fecha e devolve o foco ao botão. O tema escuro é o padrão. A preferência salva no navegador também se aplica ao login e atualiza o favicon.

O caminho de navegação começa com **MIRA** em texto. **PROJETOS** abre a tabela, inclusive na própria tela. No detalhe do container, o nome do projeto é um link para o detalhe do projeto. Os nomes de projeto mantêm sua grafia cadastrada e o nome do container aparece em maiúsculas no caminho. A página atual é indicada para leitores de tela. Os links ganham azul e sublinhado ao passar o ponteiro e o foco de teclado tem contorno visível. O atalho **Ir para o conteúdo** aparece ao receber foco.

As quatro tabelas compartilham cabeçalhos, divisórias tracejadas nas linhas e alinhamento à direita para medidas e contagens. O título do painel fica à esquerda e os filtros à direita. Até 550 px, os filtros podem ocupar a linha seguinte. Tabelas largas rolam dentro do painel. Os rodapés mostram faixa exibida, quantidade de registros filtrados, total geral quando difere e **Página X de Y**. Cada página contém até dez registros, com controles **Página anterior** e **Próxima página** desabilitados nos limites. Busca ou seleção de filtro volta à primeira página. Ícones de ação têm nome acessível e descrição ao passar o ponteiro.

Os diálogos usam fundo escurecido, largura de até 480 px e se ajustam à área visível do celular. As notificações de conclusão aparecem abaixo à direita, são anunciadas de forma discreta por leitores de tela e, no protótipo, desaparecem após quatro segundos.

## 6. Arquitetura de deploy

### Implantação inicial

- O checkout e o `docker-compose.yml` ficam em `/srv/apps/mira`. O serviço Go escuta em `0.0.0.0:6002`. Somente o serviço HTTP do Mira entra na rede externa `proxy`, com alias `mira-http`. O Compose do Mira não publica portas no host.
- O Nginx compartilhado em `/srv/infra` atende `mira.macb.tech` e resolve `mira-http:6002` dinamicamente pela rede Docker. O vhost segue o fluxo existente de ACME, emissão pelo Certbot e redirecionamento para HTTPS. Antes de publicar, conferir DNS, upstream e disponibilidade do número `6002` no inventário da VPS.
- Configurar `TRUST_PROXY` para a origem verificada do Nginx. A rede `proxy` também conecta outros aplicativos. As rotas SSE desativam o buffering e emitem sinais periódicos em intervalo inferior ao timeout de leitura do Nginx. [Documentação do Nginx](https://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_read_timeout)
- Rodar o serviço sem root, com `.env` externo ao versionamento, `restart: unless-stopped`, encerramento por SIGTERM e healthcheck HTTP do próprio Mira. Enviar logs a stdout/stderr com rotação `json-file` de 10 MB por arquivo e três arquivos. Medir CPU, RAM e disco com os aplicativos atuais ativos antes de fixar o orçamento de recursos.
- Montar `/srv/data/mira` no serviço Go, com diretório provisionado para o UID/GID do processo e `create_host_path: false`. Gerar um backup consistente do `mira.db` com o SQLite em execução, armazenar em `/srv/backups/mira/` e testar uma restauração isolada. [API de backup do SQLite](https://sqlite.org/backup.html)
- O componente de integração com Docker é o único serviço do Mira com acesso ao socket. Ele recebe chamadas da API pela rede privada do projeto. Começar pelas consultas do inventário e acrescentar leituras de métricas e logs nas etapas correspondentes. Após validação em containers de teste, liberar `start`, `stop` e `restart` para serviços explicitamente autorizados. [Segurança do Docker Engine](https://docs.docker.com/engine/security/)

### Deploy via GitHub Actions

O vínculo operacional de cada projeto define o repositório, o workflow permitido, a referência protegida e os serviços abrangidos. O Mira usa credencial de servidor com acesso aos repositórios necessários e valida o vínculo antes de acionar o workflow. O evento `workflow_dispatch` chama o mesmo procedimento de deploy por SSH e Compose usado pelo pipeline do projeto. O workflow confere a revisão aprovada e registra sucesso depois de verificar a saúde. A configuração de portas, variáveis, volumes, redes, aliases e labels vem do Compose. [API de workflows do GitHub](https://docs.github.com/en/rest/actions/workflows#create-a-workflow-dispatch-event)

O Mira persiste o ID da execução e consulta seu andamento pelo GitHub. O deploy manual e o automático compartilham a trava de deploy da VPS para o projeto. Durante um deploy iniciado pelo painel, a API bloqueia ações de container nos serviços afetados. Antes de habilitar um projeto, validar a execução manual e a concorrência em ambiente de teste.

### Migração posterior de portas internas

Executar as mudanças nesta ordem para que o destino de cada etapa esteja livre no inventário de portas internas.

| Ordem | Serviço | Upstream atual na etapa | Novo upstream |
| --- | --- | --- | --- |
| 1 | Breve Link | `breve-link-http:6001` | `breve-link-http:6003` |
| 2 | Lembrei de Maria | `lembrei-de-maria-http:6000` | `lembrei-de-maria-http:6001` |
| 3 | Mira | `mira-http:6002` | `mira-http:6000` |

Em cada etapa, alinhar a porta de escuta, o healthcheck do serviço e o upstream do Nginx. Validar a configuração antes da recarga, verificar o healthcheck e testar o domínio por HTTPS antes de avançar. Ao concluir, a porta interna `6002` volta a ficar disponível no inventário.

## 7. Contrato mínimo da API

| Método | Rota | Resultado |
| --- | --- | --- |
| POST | `/api/auth/login` | Autentica com e-mail e senha e emite a sessão |
| GET | `/api/auth/me` | Identifica a conta da sessão com nome, e-mail e perfil |
| POST | `/api/auth/refresh` | Rotaciona o refresh |
| POST | `/api/auth/logout` | Revoga a sessão |
| GET / POST | `/api/users` | Lista contas com situação e último acesso ou cria uma conta administrativa ativa |
| PATCH | `/api/users/:id` | Edita nome, e-mail e senha opcional ou altera a situação, protegendo a conta em uso |
| GET / POST | `/api/projects` | Lista ou cria projetos |
| GET / PUT / DELETE | `/api/projects/:id` | Consulta, edita ou arquiva o projeto. O `GET` inclui a disponibilidade do workflow e a revisão apta para deploy, quando houver. O `PUT` mantém o slug e também permite reativar. O `DELETE` altera a situação para **Arquivado** |
| GET | `/api/containers` | Lista containers vinculados com projeto, nome, imagem, tag, estado e uma leitura dos recursos |
| GET | `/api/containers/:name/stats` | SSE de estado, CPU, memória em uso, limite e uptime |
| POST | `/api/containers/:name/:action` | Aceita `start`, `stop` e `restart` |
| GET | `/api/containers/:name/logs` | SSE com linhas recentes e acompanhamento enquanto o container estiver ativo |
| GET | `/api/system/metrics` | Retorna domínio, IP, CPU, vCPUs, load, memória, disco e uptime do host |
| POST | `/api/projects/:id/deploy` | Aceita o SHA da revisão aprovada, solicita a execução do workflow configurado e retorna o identificador do deploy. Um projeto com deploy em andamento devolve o conflito com a referência à operação ativa |
| GET | `/api/projects/:id/deploys` | Lista resultados do mais recente ao mais antigo e identifica a operação em andamento, quando houver |
| GET | `/api/deploys/:id/logs` | SSE com etapas acumuladas e novos eventos até o resultado final |

Erros de validação de formulário identificam o campo. A API devolve dados de conta sem credenciais ou hashes. A autorização consulta a situação da conta nas operações autenticadas. Desativar uma conta revoga suas sessões de renovação e encerra seus fluxos SSE.

Os fluxos SSE usam `text/event-stream`. O cliente encerra assinaturas ao trocar de fonte ou sair da tela e tenta reconexão com intervalo crescente após queda. A interface mantém os dados já recebidos e informa a situação do fluxo. Na retomada, combina as linhas recentes do container com o buffer e as etapas registradas do deploy com o progresso já exibido, evitando duplicatas.

## 8. Dados persistidos

```text
users          (id, name, email UNIQUE, password_hash, role, active, last_sign_in_at, created_at)
refresh_tokens (id, user_id, token_hash, expires_at, revoked)
projects       (id, name, slug UNIQUE, client, domain, description, status, created_at)
deployments    (id, project_id, revision_sha, services, workflow_run_id, workflow_run_url, status, started_at, finished_at, log)
audit_log      (id, user_id, action, target, result, created_at)
```

A coluna `email` guarda o valor normalizado. O perfil inicial e único no MVP é `admin`, exibido como **Administrador**. `last_sign_in_at` começa vazio e recebe o instante de cada login bem-sucedido. A situação da conta controla a autenticação e o acesso à API.

O vínculo de containers é resolvido pelo label Docker `mira.project=<slug>`. O status do projeto descreve seu ciclo de vida. O estado Docker descreve a execução de cada container. Arquivar conserva os registros de deploy. Cada deploy usa seu próprio identificador e a persistência garante uma operação ativa por projeto.

## 9. Estrutura prevista

```text
server/
  api/
  cmd/
    api/main.go
  internal/{auth,users,projects,containers,logs,deployments,metrics,audit,docker}/
  migrations/
  go.mod
web/
  modules/{auth,users,projects,containers,deployments}/
  src/
    app/
    shared/{api,components,styles}/
    main.tsx
  index.html
  package.json
deploy/
  nginx/
Dockerfile
docker-compose.yml
.env.example
```

| Pasta | Responsabilidade |
| --- | --- |
| `server/api/` | Rotas HTTP com chi, handlers, middleware, respostas e conexões SSE. |
| `server/cmd/` | Pontos de entrada. `cmd/api/main.go` carrega a configuração, monta as dependências e inicia o servidor. |
| `server/internal/` | Regras de negócio e integrações com SQLite, o componente Docker e o GitHub Actions. A camada HTTP chama esses pacotes. |
| `server/migrations/` | Migrações do banco de dados. |
| `web/src/main.tsx` | Ponto de entrada que monta o React no HTML. |
| `web/src/app/` | Layout do console, navegação e composição das telas. |
| `web/modules/` | Telas e lógica de cada área. **Projetos** mostra também as métricas da VPS. **Containers** inclui o painel de logs. |
| `web/src/shared/` | Cliente HTTP e SSE, componentes reutilizáveis e estilos comuns, incluindo os temas. Pode ser usado por `web/src/app/` e `web/modules/`. |
| `deploy/nginx/` | Configuração do vhost para integrar ao Nginx compartilhado em `/srv/infra`. |
| `Dockerfile` e `docker-compose.yml` | Build com etapas para `web/` e `server/` e configuração dos serviços do Mira. O contexto de build é a raiz do projeto. |

## 10. Critérios de aceite

| Fluxo | Verificação no console local |
| --- | --- |
| Acesso | Mostrar e ocultar a senha. Enviar credenciais inválidas com Enter, conferir mensagem e foco e restaurar os campos. Entrar, conferir nome no cabeçalho e último acesso em **Usuários**. Sair e conferir o foco no e-mail. |
| Navegação e tema | A mira leva ao topo de **Projetos** e o atalho no menu leva à tabela. Voltar ao projeto pelo caminho do container preserva filtros. Escape fecha o menu e devolve o foco. O tema escuro é padrão e o tema escolhido permanece no login e após recarga. |
| Movimento | Conferir que as pontas azuis seguem o ponteiro nas duas paletas e que o fundo permanece estático com movimento reduzido. |
| Projetos e VPS | A amostra começa com quatro projetos e 7 de 9 containers ativos. Conferir domínio, IP, uptime, vCPUs, load e unidades. Buscar por slug, cliente, domínio e nome de container. Conferir contagem e mensagem quando a busca fica sem resultados. |
| Cadastro e ciclo do projeto | Criar com cliente e domínio vazios, conferir **Pessoal**, **Sem domínio** e o label no detalhe vazio. Tentar slug duplicado. Editar mantendo busca e slug. Arquivar, conferir preservação do histórico e dos containers e reativar pelo lápis. |
| Recursos e inventário | Combinar busca e estado, abrir pelo nome e voltar ao projeto. Conferir que outro projeto inicia com filtros limpos. Cancelar uma parada preserva os valores. Confirmá-la atualiza ativos, CPU, memória e uptime. Conferir o limite mesmo com o container parado e o estado vazio em um projeto sem vínculos. |
| Container e logs | Abrir um ativo e conferir imagem, métricas e início do fluxo. Buscar `WARN`, esperar novas linhas, rolar para cima e conferir que o acompanhamento automático para. Pausar congela a contagem e retomar recebe novas linhas. Parar o container preserva o histórico e desabilita a pausa. Trocar de container limpa o filtro. |
| Deploy | Trocar o alvo carrega imagem e tag. Iniciar, fechar, navegar e reabrir o progresso pelo mesmo projeto. Conferir bloqueio das ações sobre o alvo durante a operação, resultado único e filtros de concluídos e falhos. A tag `indisponivel` registra falha e conserva a imagem e o estado anteriores. |
| Usuários | Criar uma conta, tentar repetir seu e-mail com outra capitalização, editar com senha vazia e conferir que as credenciais continuam válidas. Usar a conta inicial para desativar a nova conta, verificar a recusa de login e reativá-la. Conferir a proteção da conta em uso, o filtro e o último acesso. |
| Paginação, teclado e celular | Com mais de dez registros, conferir faixa, total, setas e retorno à página 1 após busca ou filtro. Em 320 px, verificar rolagem das tabelas dentro dos painéis, acesso aos filtros e ações. Usar Tab, Enter e Escape e conferir o foco inicial em **Cancelar** nas confirmações. |

Para liberar o MVP conectado, repetir os fluxos aplicáveis com a API e containers de teste. Verificar persistência de projetos, usuários e deploys após recriar o serviço, restauração do SQLite, expiração e revogação de sessão, bloqueio de contas inativas, retomada de SSE, métricas reais do host, falhas do Docker e registros de auditoria. Confirmar acesso por HTTPS via Nginx e a ausência de portas publicadas pelo Compose do Mira.

No deploy conectado, verificar a recusa de projeto ou revisão sem autorização, a retomada do progresso após sair da tela, o bloqueio das ações nos serviços afetados, a concorrência com o deploy automático e o resultado único no histórico. Uma falha anterior à aplicação do Compose deve preservar os containers. Se falhar na verificação de saúde, o painel deve indicar a falha e mostrar o estado observado no inventário.

## 11. Segurança do MVP

| Controle | Condição de entrega |
| --- | --- |
| Transporte e navegador | HTTPS, HSTS, CSP e proteção contra incorporação em frames. |
| Sessão | Cookies `HttpOnly`, `Secure` e `SameSite=Strict`. Refresh rotativo, rate limiting no login e bloqueio de contas inativas. |
| Administrador | Seed inicial com senha forte fornecida pelo ambiente. Contas criadas no console seguem o mesmo armazenamento de senhas com bcrypt. A API impede a desativação da conta autenticada. |
| Docker | Socket montado apenas no componente dedicado em rede privada. API autenticada como única consumidora, com validação das operações e autorização por serviço. |
| Deploy | Workflow e repositório definidos na configuração do servidor. A credencial do GitHub fica no servidor. As credenciais SSH da VPS ficam no pipeline existente. |
| Auditoria | Registro append-only do operador, ação, alvo, instante e resultado nas operações de container, deploy e gestão de contas. |
| Interface | Campos de projeto, nomes, tags, mensagens e linhas de log são tratados como texto. |
