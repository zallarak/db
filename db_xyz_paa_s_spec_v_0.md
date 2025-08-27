# db.xyz — Postgres‑as‑a‑Service on Proxmox LXC

**Status:** Draft (v0.1, running doc)  
**Editors:** <your names>  
**Domain:** `db.xyz`  
**Stack:** Go (backend API + CLI), Web (SPA) consuming the API (vite, react-ts), Proxmox LXC data plane

---

## 0) TL;DR
We will ship a Postgres‑as‑a‑Service that provisions **hardened unprivileged LXC containers** on a Proxmox cluster. **Control plane** runs on `api.db.xyz` (Go). Users register, can join multiple orgs, create projects and **DB instances** from simple plans. **Queueing uses Postgres** (single datastore). CLI (`dbx`) and the Web Console **on `db.xyz`** are thin clients over the same API.

*First pass (MVP):* email/password auth only; create/list/delete instances; rotate credentials; minimal metrics (CPU/MEM/disk, QPS, connections). **Backups, upgrades, IP allow‑lists, and VPN/private VLAN** are *planned* but not in the first pass; the design keeps them in mind for smooth addition.

---

## 1) Goals
- **Fast, reliable provisioning** of per‑tenant Postgres in LXC on Proxmox.
- **Isolation-first** without guest shell access (DB over TCP only).
- **Self‑serve**: open registration, multi‑org membership, invitations & SSO ready.
- **Day‑2 ops** built‑in: backups (PITR), monitoring, firewall allow‑lists, upgrades.
- **Simple API** consumable by CLI & Web UI; documented via OpenAPI.
- **Predictable SKUs**: small/medium/large instance classes with quotas.

### Non‑Goals (for v1)
- Cross‑region replication; managed read replicas (consider v1.1+).
- On‑box extensions marketplace (we’ll ship a curated allowlist only).

---

## 2) Personas & Primary Use Cases
- **Indie dev / startup engineer:** needs a secure Postgres quickly, simple backups.
- **Data/app teams:** multiple environments (dev/stage/prod), IP allow‑lists, metrics.
- **Agencies / consultants:** single user in multiple customer orgs.

---

## 3) High‑Level Architecture

```
            +-------------------+           +-----------------------+
            |  Web Console @    | <-------> |       API (Go)       |
            |    db.xyz         |           |  auth, orgs, jobs,   |
            +-------------------+           |  instances            |
                    ^                         +-----------------------+
                    |                                   |
                    |  CLI (`dbx`)                      |  Postgres (control DB + job queue)
                    +-----------------------------------+----------------------------+
                                                                                |
                                                                       +---------------+
                                                                       | Provisioner   |---> Proxmox API
                                                                       | (Go worker)   |     (token)
                                                                       +---------------+
                                                                                |
                                                                    +-----------------------+
                                                                    |   Proxmox Cluster     |
                                                                    |  LXC hosts + ZFS/Ceph |
                                                                    +-----------------------+
                                                                                |
                                                            +----------------------------------+
                                                            |   Unprivileged LXC per instance  |
                                                            |  Postgres + TLS + exporter       |
                                                            +----------------------------------+
```

### Control Plane Components
- **API Service (Go, HTTP/JSON)**: auth, orgs/users/projects, instances, simple audit, minimal metrics proxy.
- **Postgres (single datastore)**: control metadata **and** reliable job queue (tables + advisory locks / `SKIP LOCKED`).
- **Provisioner Workers**: execute jobs (create/delete/rotate creds) via Proxmox API.
- **Static assets** for Web Console served from `db.xyz`.

### Data Plane
- **Proxmox cluster** with ZFS (or Ceph). Each instance = one **unprivileged LXC**.
- **Networking (v1)**: simple TCP exposure only. **Future**: private VLAN + WireGuard, and per‑instance IP allow‑lists.

-------------------+           +-----------------------+
            |  Web Console SPA  | <-------> |       API (Go)       |
            +-------------------+           |  authz, scheduler,   |
                    ^                         |  billing, events      |
                    |                         +-----------------------+
                    |  CLI (`dbx`)                 |      ^
                    +------------------------------+      |
                                                         (msg bus)
                                                           |
                                                   +---------------+
                                                   |  Provisioner  |---> Proxmox API
                                                   |  (Go worker)  |     (token)
                                                   +---------------+
                                                           |
                                               +-----------------------+
                                               |   Proxmox Cluster     |
                                               |  LXC hosts + ZFS/Ceph |
                                               +-----------------------+
                                                         |
                                         +----------------------------------+
                                         |   Unprivileged LXC per instance  |
                                         |  Postgres + TLS + exporter +     |
                                         |  pgBackRest agent                |
                                         +----------------------------------+
```

### Control Plane Components
- **API Service (Go, HTTP/JSON)**: auth, orgs/users/projects, instances, backups, metrics, audit, billing. OpenAPI‑first.
- **Scheduler**: chooses Proxmox node/storage; enforces quotas.
- **Provisioner Workers**: execute jobs (create/upgrade/backup/restore) via Proxmox API.
- **Control DB**: Postgres (separate from customer DBs) for metadata & auth.
- **Queue**: NATS or Redis streams for job orchestration (lean toward NATS).
- **Object Store**: S3/MinIO for backup repos and artifacts.
- **Observability**: Prometheus + Alertmanager, Loki for logs, Grafana for dashboards.

### Data Plane
- **Proxmox cluster** with ZFS (or Ceph) storage. Each instance = one **unprivileged LXC**.
- **Networking**: private vLAN + WireGuard mesh for customers, or public IP with IP allow‑lists via Proxmox Firewall.

---

## 4) Identity, Auth, and Org Model

### Entities
- **User**: can belong to many **Orgs**; owns API keys; can be invited.
- **Org**: billing/governance entity; contains **Projects**.
- **Project**: grouping for **Instances** and (future) **Backups**.
- **Instance**: a Postgres server in one LXC.

### Roles (simple)
- **Owner** (org‑wide), **Admin**, **Member**, **Viewer**.

### AuthN (v1)
- **Email + password** only. Note: **2FA (TOTP)** planned post‑MVP.

### AuthZ
- Minimal role checks in code (no external policy engine). Keep it simple and obvious.

---

## 5) Domains & DNS
- `db.xyz` — hosts landing page **and** Web Console (cookie auth for logged‑in state)
- `api.db.xyz` — API (Go)
- `*.cust.db.xyz` — future per‑instance FQDNs (e.g., `pg-abc123.cust.db.xyz`) when network features land

---

## 6) API (HTTP/JSON) — v1 Surface
**Design:** REST-ish over HTTPS, cookie sessions for Web, bearer tokens / API keys for CLI. Pagination `page[size]&page[token]`. Idempotency via `Idempotency-Key`.

**Base URL:** `https://api.db.xyz/v1`

### Core Resources (v1)
- `POST /auth/register`, `POST /auth/login`, `POST /auth/logout`
- `GET /users/me`
- `GET/POST /orgs`, `GET/PATCH/DELETE /orgs/{orgId}`
- `POST /orgs/{orgId}/invitations`, `POST /invitations/{token}/accept`
- `GET/POST /orgs/{orgId}/projects`
- `GET/POST /projects/{projectId}/instances`
- `POST /instances/{id}:start|stop|reboot|rotate-creds`
- `GET /instances/{id}`, `DELETE /instances/{id}`
- `GET /instances/{id}/metrics` (minimal: CPU/MEM/disk, QPS, connections)
- `GET/POST /apikeys`
- `GET /audit-logs?orgId=...` (basic events only)

**Deferred (design for later, not in v1):**
- Backups & restore endpoints
- Resize & upgrade endpoints
- Network policy endpoints (allow‑lists, VPN)

### Example: create instance (v1)
```http
POST /v1/projects/{projectId}/instances
Content-Type: application/json
Authorization: Bearer <token>

{
  "name": "prod-db",
  "plan": "pro",           // one of: nano | lite | pro | pro-heavy
  "pgVersion": 16,
  "diskGiB": 100,
  "tags": {"env": "prod"}
}
```
Response: `202 Accepted` with job id; `GET /jobs/{id}` to track.

### OpenAPI
- Spec in repo (`/api/openapi.yaml`); Console and CLI use generated clients.

---

## 7) Provisioning & Lifecycle (happy path)
1. **API** validates plan → enqueues **CreateInstance** job into Postgres queue.
2. **Provisioner** selects Proxmox node/storage (simple headroom heuristic).
3. Clone **LXC template**; assign CTID; set CPU/RAM per plan; attach ZFS dataset for `/var/lib/postgresql`.
4. Write `/etc/tenant.env` with `PG_TENANT_USER/DB/PASS`; start CT; `pg-firstboot.service` initializes Postgres (TLS, SCRAM, role/db).
5. Register instance (name, CTID, node, IP/FQDN(optional)) in control DB.
6. Metrics target registered for minimal stats. (No user‑facing allow‑lists/VPN in v1.)

### Rotate credentials
- `POST /instances/{id}:rotate-creds` → new password generated and stored; `pg_hba.conf` remains unchanged.

### Delete
- Stop CT; destroy LXC; destroy ZFS dataset; purge control records.

---

## 8) Proxmox & LXC Hardening Baseline
- **Unprivileged CTs** only.
- CT config defaults:
```
# /etc/pve/lxc/<CTID>.conf
unprivileged: 1
features: keyctl=0,nesting=0
memory: <ramMiB>
cores: <cpu>
mp0: /rpool/ct-<CTID>-pgdata,mp=/var/lib/postgresql
net0: name=eth0,bridge=vmbr0,firewall=1,ip=<cidr>,gw=<gw>
```
- **Firewall**: CT firewall enabled by ops; in v1, rules managed operationally (no user‑facing allow‑list yet). Only `5432/tcp` should be reachable per environment policy.
- **Inside CT**:
  - Postgres 16+, `password_encryption = scram-sha-256`, `listen_addresses='*'`.
  - TLS enabled (self‑signed default; plan for managed certs later).
  - `pg_hba.conf`: `hostssl ... scram-sha-256`.
  - `postgres_exporter` for minimal metrics.
- **Data layout**: separate ZFS dataset per instance → quotas/snapshots possible later.

---

## 9) Backups & Restore
**Not in first pass.**

**Compatibility considerations:**
- Keep data paths (`/var/lib/postgresql`) and instance metadata ready for pgBackRest/WAL‑G integration later.
- Reserve API shapes for future: `GET/POST /instances/{id}/backups`, `POST /instances/{id}/restore`.
- Consider storing minimal provenance (template revision, pg version) to enable future restore/upgrade logic.

---

## 10) Observability & Audit
- **Metrics (v1 minimal):** per‑instance CPU, memory, disk, QPS, connections. Simple in‑house collector/endpoint; no full Prometheus/Loki stack yet.
- **Logs:** basic API request logs with request IDs; instance logs accessible to ops only.
- **Audit log:** record mutating API calls (who/what/when). Keep schema minimal.

---

## 11) Plans (initial) & Future Quotas
**Plans (names + indicative resources; finalize later):**
- `nano` — 1 vCPU, 1–2 GiB RAM, 20–40 GiB disk, ~100 connections
- `lite` — 2 vCPU, 4 GiB RAM, 80–100 GiB disk, ~200 connections
- `pro` — 4 vCPU, 8 GiB RAM, 150–200 GiB disk, ~400 connections
- `pro-heavy` — 8 vCPU, 16 GiB RAM, 300–400 GiB disk, ~800 connections

**Quotas:** not in v1; note that org/project quotas (count, vCPU/RAM totals, backup storage) will be added later.

---

## 12) CLI (`dbx`) — UX Sketch
```
# auth & orgs
dbx login
dbx org list | select <org>

# projects & instances
dbx project create myapp

# create minimal instance
dbx instance create --project myapp \
  --name prod-db --plan pro --pg 16 --disk 100

dbx instance list --project myapp
dbx instance info prod-db

dbx instance creds rotate prod-db

dbx instance delete prod-db --force
```
- `--output json|table` and shell completion.
- No network/backup/upgrade commands in v1; will arrive with features.

---

## 13) Web Console — Key Flows
- **Hosted on `db.xyz`**. Landing page for logged‑out users; Console UI when authenticated (cookie‑based sessions).
- Signup/login → create/join org → invite members.
- Create project → create instance (plan, version, disk).
- Instance detail: status, minimal metrics, connection string, rotate credentials, delete.
- Simple activity feed (derived from audit log).

---

## 14) Security & Compliance Notes
- Secrets: store DB creds encrypted at rest (KMS/SOPS); API never logs plain creds.
- Passwords: Argon2id with tuned params; 2FA (TOTP) in v1.1.
- TLS: HSTS on console & API; perfect‑forward‑secrecy ciphers.
- Hard multi‑tenancy: one LXC per instance; no shell access exposed.
- DDoS/abuse: rate limiting per IP and per token; WAF in front of API.

---

## 15) Data Model (control plane) — Sketch
```
User(id, email, pw_hash, created_at, ...)
Org(id, name, created_at, ...)
Membership(user_id, org_id, role)
Project(id, org_id, name)
Instance(id, project_id, name, plan, pg_version, node, ctid, fqdn, status, ...)
NetworkPolicy(instance_id, exposure, allowed_cidrs[])
BackupPolicy(instance_id, schedule, retention_days)
Backup(id, instance_id, repo_url, started_at, completed_at, status, size_bytes)
ApiKey(id, user_id, org_id, hash, prefix, created_at, last_used_at)
AuditLog(id, actor_user, org_id, resource_urn, action, diff_json, ip, ts)
Job(id, type, payload_json, status, created_at, updated_at)
```

---

## 16) Proxmox Integration Details
- **Auth to Proxmox:** API token with least privileges (LXC create/config/start/stop/destroy, storage, firewall read/update).
- **Node selection:** simple heuristic by free RAM/CPU; future: binpack/spread strategies.
- **Storage:** ZFS pool `rpool` with dataset per instance `rpool/ct-<CTID>-pgdata` (quota per plan when quotas arrive).
- **Networking:** `vmbr0` bridge; CT firewall enabled by ops. **Future:** private VLAN + WireGuard, and user‑managed allow‑lists via API/UI.
- **Templates:** Debian 12 LXC with `pg-firstboot.service` baked in.

---

## 17) Versioning & Upgrades
- API: semantic versioning; breaking changes gated behind `/v2`.
- Postgres: default latest stable minor; major upgrades via blue/green flow.
- LXC Template revisions tracked; instances created from a given template record `template_rev`.

---

## 18) SLOs / Reliability
- API availability: 99.9% monthly.
- Instance RPO: 15 minutes (via WAL streaming).
- Instance RTO: < 30 minutes for restore‑to‑new.
- Support response (internal): business hours for v1.

---

## 19) Milestones (proposed)
- **M0** (Infra): Proxmox LXC template, ZFS layout, manual create → working DB.
- **M1** (Control plane skeleton): email/password auth, orgs/projects, jobs in Postgres; minimal audit.
- **M2** (Instances v1): create/list/delete, rotate creds, minimal metrics; CLI & Web Console basics.
- **M3** (Hygiene): hardening review, API docs, first customer trial.
- **M4 (Post‑MVP targets)**: backups & restore, IP allow‑lists/VPN, upgrades, quotas.

---

## 20) Open Questions / Decisions
- Minimal metrics implementation detail (collector agent vs. scraping exporter endpoints).
- Network roadmap order: IP allow‑lists first vs. private VLAN + WireGuard?
- Backup/restore design: pgBackRest vs WAL‑G; per‑instance repos vs shared.
- Plan shapes & defaults; connection limits per plan.
- Billing model & provider timeline.

---

## 21) Appendix — First‑boot Script (CT)
Skeleton of `pg-firstboot.sh` & systemd unit are available from our previous notes. Ensure:
- sets `password_encryption = scram-sha-256`, enables TLS, and creates tenant role/db.
- configures UFW and writes `pg_hba.conf` with **CIDR‑scoped `hostssl`** entries.
- registers `postgres_exporter` service.

---

## 22) Interface, miscellaneous
Use shadcn for initial interface

*End of v0.1. This doc is living; add comments inline and expand sections as we finalize choices.*

