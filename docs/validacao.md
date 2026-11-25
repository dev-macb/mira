# Validação

Abra [`prototipo.html`](../prototipo.html) num navegador e use as credenciais preenchidas. Recarregue o arquivo entre cenários quando precisar restaurar a amostra. O tema pode permanecer salvo.

| Fluxo | Verificação manual |
| --- | --- |
| Login e saída | Altere a senha e confirme a mensagem de erro e o foco no e-mail. Use **Restaurar campos**, entre, confira o último acesso em **Usuários** e saia pelo menu. |
| Projetos e navegação | Confira os quatro projetos e os sete containers ativos de nove. Busque por domínio ou nome de container. Edite um projeto e confira a busca preservada. Abra um container pelos resultados filtrados, volte ao projeto e confira os filtros. |
| Cadastro e arquivamento | Tente cadastrar um slug existente. Crie um projeto sem cliente ou domínio e confira **Pessoal**, **Sem domínio**, a orientação sobre o label e o deploy desabilitado. Arquive um projeto com histórico, confira os vínculos e reative-o pelo lápis. |
| Containers e logs | Cancele uma parada e confira o estado. Confirme a parada e confira métricas, limite de memória e o buffer visível. Inicie novamente e confira a chegada de novas linhas, a busca por `INFO`, a pausa e a retomada. |
| Deploy | Inicie um deploy no projeto **SaaS pessoal**, feche o diálogo e reabra por **Acompanhar deploy**. Confira o bloqueio das ações do alvo durante a operação e o resultado único no histórico. Repita com a tag `indisponivel` e confira a falha e a imagem anterior preservada. |
| Usuários | Crie uma conta, tente repetir o e-mail com outra capitalização e edite a conta deixando a nova senha vazia. Confira o acesso com a senha inicial. Volte à conta administrativa inicial, desative a conta criada e confira a recusa do login. |
| Interface | Crie mais de dez registros numa tabela e confira a paginação e seu reinício ao filtrar. Em 320 px, confira a rolagem das tabelas dentro dos painéis. Use Tab, Enter e Escape, confira o foco em **Cancelar**, alterne o tema e teste o movimento reduzido. |

A versão conectada tem verificações próprias de persistência, sessão, SSE, concorrência, Docker e auditoria. Consulte os [critérios de aceite](./especificacao.md#10-critérios-de-aceite) antes de liberar o MVP na VPS.

## Base da semana 1

- Siga o [guia de desenvolvimento](../README.md#desenvolvimento-local), abra a página servida pelo Go e confira `200` e `{"status":"ok"}` em `/api/saude`.
- Reinicie o Go e confirme a abertura do SQLite com as migrações já aplicadas.
- Na VPS, siga a [validação da implantação](../deploy/README.md#validação) para conferir a consulta pela rede privada e o acesso por HTTPS sem publicar portas do Mira no host.
