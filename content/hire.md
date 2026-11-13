+++
title = "hire"
description = "backend engineering, on contract. go services on postgres and gcp, from the design doc to a running environment."
template = "hire.html"
updated = 2026-10-28

[extra]
available_from = "december"
slots = 2
# leave empty to drop the calendar links
booking_url = "https://cal.com/nagasrinath/quick-chat"

intro = "i do backend engineering on contract, for founders and small teams building a real product with complicated logic underneath."
based = "i'm in hyderabad (ist, utc+5:30) and work async-first."
stack = "go · typescript · postgres · sqlc · goose · cloud tasks · pub/sub · gcp · cloud run · terraform · github actions · testcontainers"

[[extra.scope]]
title = "new backends in go"
body = "rest apis, background workers and scheduled jobs, built as a modular monolith."

[[extra.scope]]
title = "data and domain logic"
body = "postgres schema design, reversible migrations, and the tricky parts: money, pricing, tax, multi-tenant access."

[[extra.scope]]
title = "cloud, ci and integrations"
body = "gcp and cloud run via terraform, ci with lint and vulnerability checks, and adapters for payments, sms and suppliers."

[extra.work]
name = "Srilin Electronics"
year = "2026"
status = "delivered"
summary = "an india-first electronics procurement platform. it prices customer boms across digi-key, mouser, element14 and robu, with indian customs duty and gst baked into the landed cost."
stats = [
  { n = "28", label = "go modules" },
  { n = "~149k", label = "lines of go" },
  { n = "85", label = "postgres tables" },
  { n = "187", label = "test files" },
]

[extra.method]
steps = [
  { when = "first", body = "a free 15-minute intro call." },
  { when = "week 1", body = "a design doc and a fixed-price quote." },
  { when = "week 2 on", body = "the build, in small reviewed steps." },
  { when = "last week", body = "handover, with docs and a walkthrough." },
]

[extra.pricing]
note = "usd for clients outside india, inr inside india, gst added where it applies. projects are 50% upfront and 50% on handover; audits are paid upfront; retainers are billed monthly in advance."
plans = [
  { name = "backend audit", usd = "$1,200", inr = "₹60,000", body = "one-time. the fee is credited against a project if one follows." },
  { name = "project build", usd = "$3,500", inr = "₹1,75,000", from = true, body = "fixed scope, fixed price." },
  { name = "monthly retainer", usd = "$1,000", inr = "₹50,000", per = "month", body = "up to 20 hours a month, including on-call and performance tuning." },
]

[extra.fit]
good = "founders or small teams building a real product, where the hard part is the logic."
bad = [
  "when a pixel-perfect frontend or a mobile app is the main thing",
  "when it has to ship in days, with no time to design",
  "full-time hiring",
]

[[extra.faqs]]
q = "do you work with clients outside india?"
a = "yes. i work async in IST and keep a few hours of overlap with US or EU mornings for calls. international work is invoiced in USD."

[[extra.faqs]]
q = "who owns the code?"
a = "you do, from the first commit. everything lives in your github and your cloud accounts; nothing is locked to me."

[[extra.faqs]]
q = "can you join an existing codebase?"
a = "yes - start with a backend audit or a scoped project. i usually begin with a short review so we agree on what to fix first."

[[extra.faqs]]
q = "do you only work on gcp?"
a = "gcp is where i am deepest, but the patterns carry over. tell me your stack on the call and i will be straight about it."

[[extra.faqs]]
q = "do you do frontend?"
a = "lightly - api clients and internal consoles in typescript. if the ui is the main deliverable, you will want a frontend specialist."

[[extra.faqs]]
q = "how do payments work?"
a = "projects are 50% upfront and 50% on handover. audits are paid upfront, retainers monthly in advance."
+++
