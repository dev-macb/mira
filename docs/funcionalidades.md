# Funcionalidades

Abra [`prototipo.html`](../prototipo.html) no navegador e entre com os campos preenchidos. A demonstração inicia com quatro projetos e nove containers, dos quais sete estão ativos.

| Tela | Operações disponíveis no protótipo |
| --- | --- |
| Acesso ao console | Entrar com conta ativa, mostrar ou ocultar a senha e restaurar as credenciais de demonstração após erro. |
| Projetos | Consultar métricas simuladas da VPS, buscar projetos e abrir o detalhe, cadastrar e editar projetos. A busca inclui nome, slug, cliente, domínio e nome de container vinculado. |
| Detalhe do projeto | Ver os recursos agregados e o inventário, filtrar containers por nome e estado, arquivar o projeto e consultar ou filtrar o histórico de deploys. |
| Detalhe do container | Consultar imagem, estado, CPU, memória e uptime, confirmar ações de início, parada ou reinício e buscar, pausar ou retomar logs. |
| Deploy manual | Escolher um container do projeto, imagem e tag, acompanhar as etapas e consultar o resultado no histórico. |
| Usuários | Buscar, criar, editar, ativar e desativar contas administrativas e consultar o último acesso. |
| Menu da conta | Ir à tabela de projetos, abrir **Usuários**, alternar o tema e sair. A mira do cabeçalho leva ao topo da tela **Projetos**. |

As tabelas paginam de dez em dez registros. Alterar a busca ou o filtro volta à primeira página da tabela. Voltar de um container ao mesmo projeto mantém os filtros e as páginas de containers e deploys. Abrir outro projeto reinicia esses controles.

Sair encerra o acesso ao console e volta ao login. As alterações na amostra permanecem em memória até recarregar. A preferência de tema pode persistir entre recargas conforme o armazenamento local do navegador.

Os [requisitos funcionais do MVP](./especificacao.md#4-requisitos-funcionais-da-implementação) descrevem a evolução desses fluxos com autenticação no servidor, dados persistentes, operações reais no Docker, métricas do host, SSE e auditoria.
