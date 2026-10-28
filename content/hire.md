+++
title = "hire"
description = "proposal: backend engineering, on contract. go services on postgres and gcp, from the design doc to a running environment."
template = "hire.html"
updated = 2026-10-28

[extra]
heading = "proposal: backend engineering, on contract"
available_from = "december"
slots = 2
# leave empty to hide the "$ open" line in the next-step block
booking_url = "https://cal.com/nagasrinath/quick-chat"

meta = [
  { k = "to", v = "founders & small teams building a real product" },
  { k = "based", v = "hyderabad · IST (UTC+5:30) · async-first" },
]

summary_lead = "backends built right the first time."
summary = "i design and ship go services on postgres and gcp — from the design doc to a running environment, with the migrations, CI and infrastructure that keep it sane after i hand it over. you own every line from the first commit."
stack = "go · typescript · postgres · sqlc · goose · cloud tasks · pub/sub · gcp · cloud run · terraform · github actions · testcontainers"

[[extra.scope]]
title = "new backends in go"
body = "apis, background workers and scheduled jobs, built as a modular monolith you can split later by moving a package — not rewriting."
items = ["rest api with a generated client", "workers on queues and pub/sub", "auth, tenancy and audit built in"]

[[extra.scope]]
title = "data & domain logic"
body = "the schema designed before the code, and the business rules that are easy to get subtly wrong — money, pricing, tax, multi-tenant access."
items = ["postgres schema and typed queries", "reversible migrations", "money in minor units, done right"]

[[extra.scope]]
title = "cloud, ci & integrations"
body = "environments as code, a pipeline that blocks mistakes before review, and third-party apis behind adapters you can swap."
items = ["gcp / cloud run via terraform", "ci with lint, vuln and drift checks", "payment, sms and supplier adapters"]

[extra.work]
name = "srilin electronics"
meta = [
  { k = "client", v = "srilin electronics — india-first electronics procurement" },
  { k = "year", v = "2026" },
]
status = "delivered"
story = [
  { k = "problem", v = "turn a messy customer BOM into a trustworthy quote — priced across several global distributors, with indian customs duty and GST worked into the landed cost." },
  { k = "approach", v = "a go modular monolith with 28 bounded contexts, designed in a frozen HLD and schema before code. supplier and payment integrations sit behind adapters; custom linters in CI enforce tenant isolation and module boundaries." },
  { k = "result", v = "delivered the full flow end to end: BOM upload → column mapping → multi-supplier pricing → quote, priced against digi-key, mouser, element14 and robu. new suppliers and gateways plug in without touching business code." },
]
size = [
  { k = "go code", v = "~149k lines across ~509 files" },
  { k = "modules", v = "28, boundaries enforced in CI" },
  { k = "postgres", v = "85 tables, reversible migrations" },
  { k = "tests", v = "187 test files, testcontainers on real postgres" },
]
figure = "one image, five deployables on cloud run (asia-south1)"
decisions = [
  "anti-IDOR linter: every query on a tenant-owned table must carry the org check, or CI fails.",
  "landed-cost pricing: HSN-keyed customs duty, IGST, spares and attrition rules — and an explanation when a line loses.",
  "migrations run up → down → up in CI on every change; sqlc drift fails the build.",
  "fixtures mode with synthetic suppliers: runs locally with no cloud account on day one.",
]

[extra.method]
lead = "design first. ship in small, reviewed steps."
steps = [
  { when = "day 1", title = "intro call", body = "15 minutes, free. you describe the problem; i tell you honestly whether i am the right fit." },
  { when = "week 1", title = "design & quote", body = "a short design doc — modules, schema, infra — with a fixed price and milestones you sign off on." },
  { when = "weeks 2+", title = "build", body = "small, reviewed commits in your repo, deployed as they land. async updates, a weekly demo." },
  { when = "last week", title = "handover", body = "docs, runbooks and a recorded walkthrough so your team can own it — or keep me on a retainer." },
]

[extra.pricing]
lead = "fixed scopes, not hours. every engagement starts with a written scope you sign off on."
note = "usd for clients outside india · inr for clients in india, GST added where applicable. an audit fee is credited against a project that follows it."
plans = [
  { name = "backend audit", usd = "$1,200", inr = "₹60,000", basis = "one-time", items = ["database & schema review", "architecture blueprint", "written report, prioritised fixes"] },
  { name = "project build", usd = "$3,500", inr = "₹1,75,000", from = true, basis = "fixed scope", items = ["end-to-end backend", "api design & docs", "deployment, ci/cd, terraform", "handover & walkthrough"] },
  { name = "monthly retainer", usd = "$1,000", inr = "₹50,000", per = "month", basis = "up to 20 hrs", items = ["on-call backend & infra", "performance tuning", "priority support on weekdays"] },
]

[extra.fit]
good = [
  "are a founder or small team building a real product, not a demo",
  "have logic that is genuinely tricky — money, pricing, tax, multi-tenant data",
  "want to own the code and run it without me later",
]
bad = [
  "need a pixel-perfect frontend or mobile app as the main deliverable",
  "want something live by friday with no design step",
  "are looking for a full-time hire (happy to chat, though)",
]

[[extra.faqs]]
q = "do you work with clients outside india?"
a = "yes. i work async in IST and keep a few hours of overlap with US or EU mornings for calls. international work is invoiced in USD."

[[extra.faqs]]
q = "who owns the code?"
a = "you do, from the first commit. everything lives in your github and your cloud accounts; nothing is locked to me."

[[extra.faqs]]
q = "can you join an existing codebase?"
a = "yes — start with a backend audit or a scoped project. i usually begin with a short review so we agree on what to fix first."

[[extra.faqs]]
q = "do you only work on gcp?"
a = "gcp is where i am deepest, but the patterns carry over. tell me your stack on the call and i will be straight about it."

[[extra.faqs]]
q = "do you do frontend?"
a = "lightly — api clients and internal consoles in typescript. if the ui is the main deliverable, you will want a frontend specialist."

[[extra.faqs]]
q = "how do payments work?"
a = "projects are 50% upfront and 50% on handover. audits are paid upfront, retainers monthly in advance."

[extra.next]
lead = "tell me what you are building — a few lines is enough. if i'm not the right person, i'll say so and point you to someone who is."
+++
