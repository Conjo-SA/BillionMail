# Deploy do BillionMail (perfil "somente ENVIO via API") ao lado do CapRover

Este guia descreve como subir o BillionMail **apenas para disparo de emails via API**,
rodando **em paralelo** ao CapRover na mesma VPS Ubuntu — não dentro do CapRover.

## Por que ao lado do CapRover e não dentro dele?

O BillionMail é um *appliance*: o container `core` orquestra os outros containers via
`/var/run/docker.sock` e roda `fail2ban` com `NET_ADMIN`/`NET_RAW`. Esse modelo não cabe
no Docker Swarm (base do CapRover), que trabalha com "1 app = 1 serviço".

A solução é rodar o `docker-compose` nativo do BillionMail no mesmo host. Não há conflito
de portas porque:

| | CapRover | BillionMail |
|---|---|---|
| HTTP/HTTPS | **80 / 443** | painel em **8080 / 8443** |
| Email | — | **25** (saída + bounces) |

Arquivos deste perfil:
- `docker-compose.billionmail.yml` — stack enxuto: `postgres, redis, rspamd, postfix, core`
- `.env.billionmail.example` — modelo de variáveis

---

## Fase 0 — Pré-requisitos de entregabilidade (faça ANTES)

Sem isto, o email vai para spam ou nem sai.

1. **Porta 25 de saída liberada.** Muitos provedores (Hetzner, DO, OVH, AWS…) bloqueiam.
   Teste na VPS:
   ```bash
   nc -zv gmail-smtp-in.l.google.com 25
   ```
   Se falhar (timeout), abra um ticket no provedor pedindo a liberação da porta 25 de saída.

2. **rDNS / PTR.** No painel da VPS, configure o PTR do IP público para `mail.seudominio.com`
   (o mesmo valor de `BILLIONMAIL_HOSTNAME`). Confira:
   ```bash
   dig -x <IP_DA_VPS> +short
   ```

> ⚠️ O IP é compartilhado com seus apps do CapRover. A reputação de envio é única para o
> servidor inteiro. Mantenha volume e conteúdo saudáveis.

---

## Fase 1 — DNS (no provedor do seu domínio)

Substitua `seudominio.com` e `<IP_DA_VPS>`:

| Tipo | Nome | Valor |
|---|---|---|
| A | `mail.seudominio.com` | `<IP_DA_VPS>` |
| TXT (SPF) | `seudominio.com` | `v=spf1 a mx ip4:<IP_DA_VPS> ~all` |
| TXT (DMARC) | `_dmarc.seudominio.com` | `v=DMARC1; p=none; rua=mailto:dmarc@seudominio.com` |
| MX | `seudominio.com` | `mail.seudominio.com` (prioridade 10) |

O **DKIM** é gerado pelo painel do BillionMail na Fase 4 — você volta aqui e adiciona o TXT.

---

## Fase 2 — Deploy na VPS

```bash
# 1. Clonar o repositório (na pasta que preferir)
git clone <url-do-repo> billionmail && cd billionmail

# 2. Criar e editar o .env
cp .env.billionmail.example .env
nano .env
#   - BILLIONMAIL_HOSTNAME=mail.seudominio.com
#   - SafePath= (algo difícil de adivinhar)
#   - ADMIN_PASSWORD= (troca depois na UI também)
#   - confirme TZ e, se 172.66.1.0/24 colidir, troque IPV4_NETWORK

# 3. Subir o stack enxuto
docker compose -f docker-compose.billionmail.yml up -d

# 4. Acompanhar
docker compose -f docker-compose.billionmail.yml ps
docker compose -f docker-compose.billionmail.yml logs -f core-billionmail
```

### Firewall (ufw)

```bash
sudo ufw allow 25/tcp       # email (saída de fato usa conexões de saída; 25 entra p/ bounces)
sudo ufw allow 8443/tcp     # painel admin (HTTPS)
# 80/443 continuam do CapRover; não mexa neles
```

Se usar firewall do provedor (cloud firewall), libere os mesmos.

---

## Fase 3 — Primeiro acesso ao painel

Abra:

```
https://<IP_DA_VPS>:8443/<SafePath>
```

(O `<SafePath>` é o valor que você pôs no `.env`.) O certificado será autoassinado no início —
aceite o aviso do navegador. Faça login com `ADMIN_USERNAME` / `ADMIN_PASSWORD` e **troque a senha**.

---

## Fase 4 — Adicionar domínio e gerar DKIM

1. No painel: **Mail Domains / Domains → Add Domain** → `seudominio.com`.
2. O painel gera o registro **DKIM**. Copie o TXT mostrado (algo como
   `dkim._domainkey.seudominio.com  =  v=DKIM1; k=rsa; p=...`).
3. Vá ao DNS do domínio e crie esse TXT.
4. Volte ao painel e clique em **verificar** — SPF, DKIM e DMARC devem ficar verdes.

---

## Fase 5 — Criar template + API Key

1. **Email Templates** → crie um template (assunto + corpo HTML). Pode usar variáveis de
   `attribs`, ex. `{{nome}}`.
2. **Mail API** (Batch Mail → API) → **criar**:
   - associe ao template e a um grupo de contatos
   - defina o remetente `addresser` (ex. `noreply@seudominio.com`)
   - copie a **API Key** gerada.

---

## Fase 6 — Enviar via API (integração)

Envio individual:

```bash
curl -X POST "https://<IP_DA_VPS>:8443/batch_mail/api/send" \
  -H "x-api-key: SUA_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "addresser": "noreply@seudominio.com",
    "recipient": "destino@exemplo.com",
    "attribs": { "nome": "Fulano" }
  }'
```

Envio em lote:

```bash
curl -X POST "https://<IP_DA_VPS>:8443/batch_mail/api/batch_send" \
  -H "x-api-key: SUA_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "addresser": "noreply@seudominio.com",
    "recipients": ["a@exemplo.com", "b@exemplo.com"],
    "attribs": { "nome": "Cliente" }
  }'
```

> Como o painel está em porta dedicada (8443) com cert autoassinado, em chamadas server-to-server
> use o IP/porta e, se necessário, ignore a validação TLS (`curl -k`) **apenas** enquanto não
> houver um certificado válido. Para produção, instale um cert válido no domínio do painel.

---

## Fase 7 — Testar entregabilidade

1. Envie um email para o endereço que o site **https://www.mail-tester.com** fornece.
2. Veja a nota (mira 9–10/10). Ele aponta SPF/DKIM/DMARC/rDNS faltando.
3. Logs úteis:
   ```bash
   docker compose -f docker-compose.billionmail.yml logs -f postfix-billionmail
   docker compose -f docker-compose.billionmail.yml logs -f rspamd-billionmail
   ```
   O rspamd é quem assina o DKIM na saída.

---

## Operação

```bash
# parar / iniciar / atualizar imagens
docker compose -f docker-compose.billionmail.yml down
docker compose -f docker-compose.billionmail.yml up -d
docker compose -f docker-compose.billionmail.yml pull && \
  docker compose -f docker-compose.billionmail.yml up -d
```

Dados persistem em volumes-bind na pasta do projeto (`postgresql-data/`, `redis-data/`,
`rspamd-data/`, `postfix-data/`, `logs/`). Inclua-os no seu backup.

---

## Quando quiser recebimento/webmail/SMTP autenticado

Este perfil corta `dovecot` e `roundcube`. Para reativar (IMAP/POP, webmail, submission
autenticada em 587/465), use o `docker-compose.yml` completo do projeto e exponha as portas
de email correspondentes — lembrando de remapear o painel para 8080/8443 para não colidir
com o CapRover.
