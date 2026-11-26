# Plano de implementação do Mira

O protótipo navegável serve de referência para a interface. A primeira versão conectada usa a infraestrutura compartilhada da VPS e entrega login, projetos e inventário de containers em modo de consulta. As ações de container e o deploy via GitHub Actions entram após essa validação. [Protótipo](./prototipo.html) · [Especificação do MVP](./docs/especificacao.md)

## Premissas do cronograma

A estimativa é de **oito semanas de trabalho para uma pessoa desenvolvedora**, contadas a partir da disponibilidade da VPS, do acesso ao DNS e de containers de teste. As semanas indicam a ordem e a duração estimadas das entregas. As datas de calendário dependem da confirmação do início.

## Cronograma

| Semana | Entrega | Condição para concluir |
| --- | --- | --- |
| 1 | **Base do projeto e integração com a VPS.** Entregar React e API no mesmo serviço Go, configurar SQLite e migrações, montar o componente Docker em modo de leitura e preparar `mira.macb.tech` no Nginx/Certbot compartilhado. | O Mira responde por HTTPS via `mira-http:6002`, sem porta publicada no host. A API lê o inventário pela rede privada. |
| 2 | **Autenticação.** Preparar o administrador inicial e ligar a tela de login à sessão, à renovação, à saída e ao registro de último acesso. | Login, renovação e logout funcionam com sessão persistida. Contas inativas têm acesso bloqueado. **RF-01** |
| 3 | **Projetos e inventário.** Ligar cadastro, consulta, edição, arquivamento e reativação ao banco. Resolver os containers pelo label do projeto. | A primeira versão conectada permite entrar, cadastrar um projeto e consultar os containers vinculados. Projetos persistem após recarregar. **RF-02 e RF-03 (inventário)** |
| 4 | **Estado, métricas e ações de container.** Ligar as estatísticas ao console e implementar início, parada e reinício com confirmação, autorização por serviço e auditoria. | Em containers de teste autorizados, a interface reflete o estado confirmado e os agregados. A auditoria registra o resultado das ações. **RF-03 e RF-04** |
| 5 | **Logs e métricas da VPS.** Entregar o fluxo de logs com pausa e reconexão e substituir as métricas simuladas do host por leituras reais. | O console exibe dados da VPS. Os logs preservam o buffer após reconectar e encerram a conexão ao sair do detalhe. **RF-05 e RF-06** |
| 6 | **Deploy via GitHub Actions e histórico.** Integrar o acionamento manual do workflow autorizado de um projeto de teste e reaproveitar o deploy via Compose. Persistir a execução e bloquear ações concorrentes nos serviços afetados. | O progresso pode ser retomado após sair da tela. A execução valida a revisão, verifica a saúde após aplicar o Compose e grava um único resultado. **RF-07 e RF-08** |
| 7 | **Gestão de usuários.** Entregar a criação, edição, ativação e desativação de contas, com proteção da conta em uso, revogação de sessões e auditoria. | Os fluxos de conta funcionam de ponta a ponta. Conta desativada perde acesso e as mudanças ficam registradas. **RF-09** |
| 8 | **Validação e publicação.** Executar os critérios de aceite, testar a restauração do SQLite e medir o consumo do Mira junto aos demais serviços. | O MVP é liberado depois de validar os fluxos funcionais, o backup e os controles de segurança aplicáveis. |

## Dependências e marcos

**Acompanhamento:** a etapa 1 está publicada na VPS. A implementação local da etapa 2 já inclui a conta inicial, login, sessão rotativa e tela de acesso; a publicação desta versão requer configurar `MIRA_JWT_SECRET`, `MIRA_ADMIN_EMAIL` e `MIRA_ADMIN_PASSWORD` na VPS. A validação operacional de renovação, logout e conta inativa continua pendente.

- **Início da semana 1** · A base local inclui React e API no mesmo serviço Go, SQLite com migrações e o componente de inventário em rede privada. O [roteiro de implantação](./deploy/README.md) prepara o Compose e os vhosts. A conclusão da semana depende da publicação e da validação com DNS, HTTPS e containers de teste na VPS.
- **Fim da semana 3** · Validar a primeira entrega navegável com autenticação, projetos e inventário em consulta. Confirmar que as operações de escrita estão indisponíveis nessa etapa.
- **Fim da semana 6** · Validar ações Docker, retomada dos fluxos SSE e acompanhamento do workflow em um projeto de teste. Exercitar falhas do workflow, a concorrência entre deploy manual e automático e o bloqueio de ações durante deploys solicitados pelo painel.
- **Fim da semana 8** · Confirmar persistência, revogação de sessões, auditoria e acesso somente por HTTPS antes de liberar o MVP. [Critérios de aceite](./docs/especificacao.md#10-critérios-de-aceite) · [Controles de segurança](./docs/especificacao.md#11-segurança-do-mvp)
- **Etapa posterior à implantação inicial** · Reorganizar as portas internas na ordem Breve Link `6001 → 6003`, Lembrei de Maria `6000 → 6001`, Mira `6002 → 6000`. Validar cada serviço antes de avançar. [Migração de portas internas](./docs/especificacao.md#migração-posterior-de-portas-internas)

O acompanhamento usa os identificadores RF-01 a RF-09 da especificação. As funcionalidades posteriores ao MVP ficam no backlog descrito no escopo. [Requisitos funcionais](./docs/especificacao.md#4-requisitos-funcionais-da-implementação) · [Escopo](./docs/especificacao.md#1-objetivo-e-escopo)
