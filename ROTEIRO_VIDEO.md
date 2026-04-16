# 🎬 ROTEIRO DO VÍDEO — Teste Korp (Vitor)

> Use este documento como guia para gravar seu vídeo de apresentação. Cada seção é uma **cena** que você vai seguir com a tela compartilhada.

---

## CENA 1 — Introdução (30 segundos)

**[Tela: tela inicial do sistema no localhost:4200]**

> "Olá pessoal da Korp! Sou o Vitor e esta é a minha solução para o desafio técnico de Emissão de Notas Fiscais."
>
> "A aplicação foi desenvolvida com **Angular 17** no frontend usando Standalone Components e um design visual em CSS puro com a tendência de **Glassmorphism**. No backend, estruturei dois microsserviços independentes em **Golang** usando o framework **Gin** e o ORM **GORM** para persistência em banco de dados real."

---

## CENA 2 — Cadastro de Produtos + IA (1 minuto)

**[Tela: Navegar até a aba "Em Estoque"]**

> "O primeiro dos meus serviços é o de **Estoque**, rodando na porta 8081. Ele controla produtos, descrições e saldos, e expõe uma API REST completa."

**[Ação: Digitar no campo 'Código': `GABINETE-ATX`]**

> "Para o cadastro, temos os campos obrigatórios: Código, Descrição e Saldo. Mas aqui eu implementei um **requisito opcional de Inteligência Artificial**."

**[Ação: Clicar no botão roxo '✨ IA' — esperar o loader girar — a descrição vai aparecer sozinha no campo abaixo]**

> "Ao invés do usuário digitar a descrição manualmente, ele pode clicar nesse botão. O Angular faz uma chamada HTTP ao meu endpoint `POST /api/products/ai/generate` no backend Go. O texto é gerado automaticamente com contexto corporativo, simulando a integração com um modelo de linguagem."
>
> "Isso economiza tempo do operador e demonstra uma aplicação real de IA no sistema."

**[Ação: Colocar saldo 10 e clicar em "Cadastrar"]**

> "Pronto. Produto persistido com sucesso no banco de dados SQLite do microsserviço de Estoque."

**[Ação: Cadastrar mais 1 produto rapidamente — ex: `MOUSE-GAMER`, IA, saldo 5]**

---

## CENA 3 — Abertura de Nota Fiscal (1 minuto)

**[Tela: Navegar até a aba "Faturamento"]**

> "Agora no segundo microsserviço: o de **Faturamento**, porta 8082. Ele possui seu próprio banco de dados isolado, cumprindo a arquitetura de microsserviços."

**[Ação: No dropdown, selecionar o produto "GABINETE-ATX" e digitar quantidade 3]**

> "A nota puxa os produtos do serviço de Estoque via HTTP e usa **RxJS BehaviorSubject** como Single Source of Truth no Angular. Isso significa que os dados são reativos — qualquer atualização reflete instantaneamente em todas as telas."

**[Ação: Clicar em '+ Adicionar Produto', selecionar "MOUSE-GAMER", quantidade 2]**

> "A nota fiscal suporta múltiplos produtos com suas respectivas quantidades, exatamente como o requisito pedia. Usei **FormArray** do ReactiveFormsModule para gerenciar os itens dinamicamente."

**[Ação: Clicar em "Abrir Nota Fiscal" — esperar o toast verde]**

> "A nota foi aberta com numeração sequencial automática e status 'Aberta'. Todos os dados estão persistidos."

---

## CENA 4 — Impressão de Nota Fiscal (1 minuto)

**[Tela: Ainda no Faturamento, olhando a nota que acabamos de criar]**

> "Agora o fluxo mais importante: a **Impressão**."

**[Ação: Clicar no botão verde '🖨️ Imprimir Nota']**

> "Ao clicar em Imprimir, o Angular mostra um indicador de processamento. No backend, o Faturamento inicia uma **transação ACID** e chama o microsserviço de Estoque via HTTP para abater os saldos."
>
> "Esse abatimento é **atômico** — a query SQL usa `UPDATE balance = balance - X WHERE balance >= X`. Isso previne saldo negativo e trata completamente o requisito opcional de **Concorrência**. Se duas notas tentarem usar o último item simultaneamente, apenas uma vai conseguir."

**[O toast verde aparece, status muda para "Fechada", e um MODAL com a Nota Fiscal formatada aparece na tela]**

> "Reparem que ao confirmar a impressão, o sistema abre um **modal com a Nota Fiscal Eletrônica** formatada — com logotipo KORP, série, data, e a tabela de itens. O usuário pode clicar em 'Imprimir / Salvar PDF' para gerar um documento real usando o `window.print()` do navegador."
>
> "E reparem: notas com status 'Fechada' não podem ser impressas novamente. Isso é **Idempotência** — outro requisito opcional que implementei."

**[Ação: Fechar o modal clicando em 'Fechar']**

**[Ação: Navegar até "Em Estoque" e mostrar os saldos atualizados]**

> "Voltando ao Estoque: o saldo do GABINETE-ATX saiu de 10 para 7 — exatamente as 3 unidades da nota. O MOUSE-GAMER saiu de 5 para 3. Perfeito!"

---

## CENA 5 — Tolerância a Falhas / Resiliência (1 minuto) ⭐ MAIS IMPORTANTE

**[Tela: Voltar ao Faturamento. Criar uma nova nota qualquer com algum produto]**

> "Agora vou demonstrar o requisito **obrigatório** de tratamento de falhas entre microsserviços."

**[Ação: Na nova nota (status Aberta), clicar no botão amarelo '⚠️ Testar Falha']**

> "Esse botão injeta um Header HTTP customizado `X-Simulate-Failure: true` que faz o Faturamento redirecionar a chamada para uma porta inexistente (9999), simulando que o serviço de Estoque caiu."

**[Ação: Apontar para o terminal do invoice-service enquanto espera — lá vai aparecer as mensagens de retry]**

> "Olhem o terminal do Go. O meu client HTTP está executando um **Pattern de Retry com Backoff**. Ele tenta 3 vezes contatar o Estoque, com intervalos de 2 segundos entre cada tentativa. Isso dá tempo para o serviço se recuperar caso seja uma falha momentânea."

**[Ação: Após ~6 segundos, o toast VERMELHO aparece na tela do Angular]**

> "Após esgotar as tentativas, o backend não quebra — ele devolve um erro tratado `503 Service Unavailable`. O Angular captura no operador `catchError()` do RxJS e exibe um feedback amigável ao usuário: 'O serviço de Estoque está indisponível'. A nota continua aberta e pode ser impressa depois que o serviço voltar."

---

## CENA 6 — Detalhamento Técnico (1 minuto)

**[Tela: Pode mostrar trechos de código no editor ou ficar na tela do sistema]**

> "Para fechar com o detalhamento técnico que vocês pediram:"
>
> "**Ciclos de vida Angular utilizados:** `ngOnInit` para carregar dados e iniciar subscriptions RxJS, e `ngOnDestroy` com `Subject/takeUntil` para desinscrever observables e evitar memory leaks."
>
> "**RxJS:** Usei `BehaviorSubject` como store centralizado no ApiService, operadores `pipe`, `tap`, `catchError` para tratamento reativo de respostas e erros, e `takeUntil` para lifecycle cleanup."
>
> "**Bibliotecas visuais:** Nenhuma biblioteca externa de UI. Todo o design foi feito com CSS puro usando design system próprio em Glassmorphism com variáveis CSS e Google Fonts Outfit. Isso deixa o bundle leve e prova domínio de CSS."
>
> "**Golang:** Gerenciamento via Go Modules independentes por microsserviço. Framework **Gin** para roteamento HTTP e **GORM** como ORM com driver SQLite para persistência real. Erros e exceções são tratados com transações ACID, `defer tx.Rollback()`, e respostas HTTP semânticas — 404 para não encontrado, 409 para saldo insuficiente, 503 para indisponibilidade de serviço."
>
> "Obrigado pela oportunidade! O repositório está público no GitHub: `Korp_Teste_Vitor`. Qualquer dúvida, estou à disposição."

---

## 📋 Checklist antes de gravar

- [ ] Microsserviço de Estoque rodando (`cd stock-service && go run main.go`)
- [ ] Microsserviço de Faturamento rodando (`cd invoice-service && go run main.go`)
- [ ] Frontend Angular rodando (`cd frontend && npm run start`)
- [ ] Banco de dados do estoque limpo (se quiser demo fresh: delete `korp_estoque.db`)
- [ ] Banco do faturamento limpo (delete `korp_faturamento.db`)
- [ ] Terminais do Go visíveis para mostrar os logs de retry na Cena 5
- [ ] Software de gravação de tela pronto (OBS, ShareX, etc.)

**Tempo total estimado do vídeo: ~5 minutos**
