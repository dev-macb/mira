# Regras de negócio

## Acesso e usuários

- Todas as contas têm perfil **Administrador**. Somente contas ativas podem entrar. O login bem-sucedido atualiza **Último acesso**.
- O nome é obrigatório e tem até 100 caracteres. O e-mail precisa ser válido, tem até 254 caracteres e é armazenado sem espaços nas extremidades e em minúsculas. A unicidade usa o e-mail normalizado.
- Uma nova conta começa ativa e exige senha. Ao editar, deixar a nova senha em branco preserva a atual. A conta em uso pode editar os próprios dados e tem a ação de desativação bloqueada.
- No MVP, a API aplica os controles de sessão. Desativar uma conta revoga seus tokens de renovação e encerra seus fluxos SSE.

## Projetos e vínculos

- Nome e slug são obrigatórios. O nome tem até 80 caracteres. O slug tem até 63 caracteres, é único e fica fixo após a criação. Ele aceita letras minúsculas e números, com hífens simples entre segmentos.
- Cliente tem até 80 caracteres e assume **Pessoal** quando vazio. Domínio e descrição são opcionais e aceitam até 253 e 2.000 caracteres, respectivamente. A tabela mostra **Sem domínio** quando o campo está vazio.
- A situação inicial é **Em desenvolvimento**. Também existem **Em produção** e **Arquivado**. Arquivar um projeto existente exige confirmação no detalhe. A edição permite reativá-lo.
- Arquivar preserva vínculos, estados, histórico e acesso às operações disponíveis. O protótipo usa o `slug` do container para simular o label Docker `mira.project=<slug>`. Um projeto sem containers vinculados tem o deploy desabilitado.

## Containers e métricas

| Estado | Conta como ativo | Ações no console |
| --- | --- | --- |
| `running`, `unhealthy` | Sim | **Reiniciar** e **Parar** |
| `stopped`, `exited` | Não | **Iniciar** |

- As ações exigem confirmação. No protótipo, iniciar ou reiniciar coloca o container em `running` e reinicia o uptime. Parar coloca o container em `stopped`. As métricas e o inventário refletem o novo estado.
- No MVP, o vínculo pelo label permite consultar o container. As ações também exigem autorização explícita para o serviço do Compose correspondente.
- Os totais do projeto somam CPU e memória em uso dos ativos. O limite agregado inclui todos os containers vinculados. O maior uptime considera apenas ativos.
- Com containers vinculados e todos inativos, o projeto exibe uso zero de CPU e memória e uptime `—`. Sem vínculos, a contagem fica em zero e as demais leituras mostram `—`. No detalhe de um container inativo, o limite de memória continua visível.

## Logs

- Cada container guarda até 200 linhas recentes, com horário local em 24 horas, nível `INFO`, `WARN` ou `ERR` e mensagem. A busca compara trechos do nível e da mensagem sem distinção de maiúsculas e minúsculas.
- Abrir o detalhe limpa a busca e posiciona a leitura no fim. Novas linhas acompanham a rolagem apenas se ela já estiver perto do fim. Pausar congela a leitura e retomar restabelece o fluxo. Containers inativos mantêm o buffer disponível para consulta.
- No MVP, pausar ou sair do detalhe fecha a conexão SSE. Ao retomar ou reconectar, a interface preserva as linhas recebidas e evita duplicatas ao combinar o trecho recente com o buffer.

## Deploy manual

- No protótipo, o alvo precisa pertencer ao projeto. A imagem é obrigatória, sem espaços e com até 180 caracteres. A tag tem de 1 a 128 caracteres. Ela começa com uma letra ASCII, um número ou `_`. As posições seguintes também aceitam `.` e `-`.
- No protótipo, o fluxo simula download, parada, recriação, início e verificação de saúde. Sucesso atualiza imagem e tag, deixa o container em `running` e reinicia o uptime. A tag `indisponivel` simula falha no download e preserva o container anterior.
- No MVP, o projeto precisa ter um workflow habilitado. O Mira valida a revisão aprovada e o workflow executa o deploy e verifica a saúde conforme o Compose do projeto.
- Cada projeto aceita um deploy em andamento. Enquanto ele roda, as ações de início, parada e reinício ficam bloqueadas para os alvos afetados. Fechar o diálogo ou mudar de tela mantém a operação em curso.
- A conclusão gera uma única entrada no topo do histórico do projeto de origem. No MVP, a entrada identifica a revisão e a execução do GitHub Actions.
