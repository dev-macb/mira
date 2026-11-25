# Guia para agentes

O Mira está na etapa de protótipo. A implementação atual está em [`prototipo.html`](./prototipo.html). Os requisitos do MVP conectado ficam na [especificação](./docs/especificacao.md). As etapas e o cronograma ficam no [plano de implementação](./PLANO.md).

## Ao trabalhar neste diretório

- Consulte a [documentação](./docs/README.md) da área afetada e confira o comportamento no HTML antes de alterar um fluxo.
- Mantenha o protótipo executável ao abrir o HTML diretamente no navegador. Use dados fictícios e preserve o isolamento das operações simuladas.
- Escreva a interface e a documentação em português do Brasil. Preserve os termos usados no console e os identificadores técnicos, como `mira.project=<slug>`.
- Reutilize os componentes, tokens e funções existentes. Siga os [padrões de código](./docs/padroes-de-codigo.md) para renderização, validação, navegação e foco.
- Ao mudar um fluxo, confira as telas que compartilham o mesmo estado, incluindo métricas, filtros, logs e histórico quando aplicável. Use o [roteiro de validação](./docs/validacao.md).
- Atualize a documentação afetada junto com a mudança. Quando uma parte do MVP for implementada, registre o que passou a funcionar e os requisitos ainda pendentes.
