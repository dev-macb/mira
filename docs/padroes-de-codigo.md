# Padrões de código

## Protótipo HTML

- Preserve a estrutura do [`prototipo.html`](../prototipo.html). O HTML define as telas e os diálogos, o CSS mora no `<style>` e a lógica principal fica no escopo de `startMira()`. Mantenha o arquivo navegável sem instalação ou build.
- Use os tokens CSS definidos em `:root` e `:root[data-theme="light"]` para cores e estados. Reutilize os estilos existentes de controles, tabelas, métricas e diálogos.
- Reaproveite funções compartilhadas nas telas. Use `active()` para classificar containers, `members()` para resolver vínculos, `pageRows()` para paginar, `renderMetrics()` para cartões e `confirmAction()` para confirmações.
- Use `textContent` para inserir texto no DOM. Ao interpolar dados em `innerHTML`, aplique `html()` aos valores dinâmicos tratados como texto ou usados em atributos HTML. Para montar seletores com identificadores dinâmicos, use `CSS.escape()` onde necessário.
- Combine restrições nativas dos formulários com `setCustomValidity()` e `reportValidity()` nas validações específicas. Limpe o erro quando o campo afetado for editado.
- Mantenha a navegação em `navigate()` e atualize a interface a partir do estado alterado. Encerre temporizadores de leitura ao trocar a fonte dos logs ou sair da tela.

## Interface e acessibilidade

- Use `pt-BR` no texto da interface e nas formatações exibidas. Preserve os valores técnicos de integração, como `running`, `unhealthy` e `mira.project`.
- Prefira os controles nativos para formulários, ações e diálogos. Associe rótulos aos campos e dê nome acessível a botões com ícone. Preserve o foco visível e o suporte a Tab, Enter e Escape.
- Direcione o foco para o título ao navegar e para **Cancelar** nas confirmações. Ao atualizar uma tabela, use a ação da linha para recuperar o foco quando ela continuar visível, com a busca como destino alternativo.
- Mantenha a rolagem horizontal das tabelas dentro dos painéis em telas estreitas. Confira as paletas clara e escura e respeite `prefers-reduced-motion` no fundo animado.

## MVP conectado

- Use `server/api/` para rotas, middleware, handlers e SSE. Mantenha regras de negócio e integrações em `server/internal/` e mudanças de esquema em `server/migrations/`.
- Em `web/`, organize a composição das telas em `app/`, os fluxos por área em `modules/` e o cliente de API, componentes reutilizáveis e estilos em `shared/`.
- Valide autorização, estado, vínculo e concorrência na API antes de executar operações no Docker. A interface informa o estado e associa erros de validação ao campo correspondente.
- Ao integrar SSE, feche assinaturas quando a tela mudar de fonte e preserve os dados já recebidos na reconexão. Mantenha os controles descritos nos [requisitos de segurança](./especificacao.md#11-segurança-do-mvp).
