# Template obrigatório do campo `desc` — kanban-force

Toda criação de card pelo MCP `kanban-force` exige que o campo `desc` siga este template. Markdown é apresentado ao usuário para revisão; **o que é enviado ao MCP deve estar em HTML**.

## Estrutura em Markdown (revisão pelo usuário)

```md
## 📑 [ID-000] Título da Task

### 🔍 1. Por quê? (Motivação)
*Descreva a origem do card: se é a correção de um bug, uma nova necessidade de negócio ou um débito técnico.*

* **Contexto:**
* **Valor/Impacto:**
* **Evidência (se bug):**

---

### 🛠️ 2. Como? (Execução)
*Descreva o plano de ação, as especificações técnicas e os critérios de aceitação.*

* **O que fazer:**
* **Especificações:**
* **Critérios de Conclusão:**
    - [ ]
    - [ ]
```

## Equivalente HTML (enviado ao MCP)

```html
<h2>📑 [ID-000] Título da Task</h2>
<h3>🔍 1. Por quê? (Motivação)</h3>
<p><em>Descreva a origem do card: se é a correção de um bug, uma nova necessidade de negócio ou um débito técnico.</em></p>
<ul>
  <li><strong>Contexto:</strong> </li>
  <li><strong>Valor/Impacto:</strong> </li>
  <li><strong>Evidência (se bug):</strong> </li>
</ul>
<hr/>
<h3>🛠️ 2. Como? (Execução)</h3>
<p><em>Descreva o plano de ação, as especificações técnicas e os critérios de aceitação.</em></p>
<ul>
  <li><strong>O que fazer:</strong> </li>
  <li><strong>Especificações:</strong> </li>
  <li><strong>Critérios de Conclusão:</strong>
    <ul>
      <li><input type="checkbox"/> </li>
      <li><input type="checkbox"/> </li>
    </ul>
  </li>
</ul>
```

## Regras

- **Markdown apenas para exibição** ao usuário antes de confirmar — nunca para envio ao MCP
- **HTML obrigatório** no campo `desc` em qualquer criação ou atualização de card
- Manter os emojis 📑, 🔍, 🛠️ exatamente como acima
- O título `[ID-000]` deve ser substituído pelo ID real da task fornecido pelo usuário
