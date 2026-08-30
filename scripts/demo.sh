#!/usr/bin/env bash
# Walks one claim through all four loops, pausing between steps so the run can be
# narrated on camera.
#
#   ./scripts/demo.sh             pauses between steps
#   ./scripts/demo.sh --no-pause  runs straight through

set -euo pipefail
BASE=${BASE:-http://localhost:8000/v1}
PAUSE=1
[ "${1:-}" = "--no-pause" ] && PAUSE=0

b(){ printf '\n\033[1;36m%s\033[0m\n' "$*"; }
step(){ if [ "$PAUSE" -eq 1 ]; then printf '\033[2m  press enter\033[0m'; read -r _; fi; }

b "1 . Reset - one motor claim, damage photo only, police report missing"
curl -s -X POST "$BASE/demo/reset" | python3 -c '
import json,sys
c = json.load(sys.stdin)["data"]
p = c["policy"]
print("  %s  stage=%s  completeness=%s" % (c["claim_number"], c["stage"], c["document_completeness"]))
print("  estimated loss %s   policy %s" % (c["estimated_loss"], p["policy_number"]))
print("  cover %s   deductible %s" % (p["max_coverage_amount"], p["deductible_amount"]))
print("  documents:", [d["document_type"] for d in c.get("documents") or []])'
step

b "2 . Upload the police report - fires Loops 1, 2 and 3 in one request"
curl -s -X POST "$BASE/claim/1/documents" -H 'Content-Type: application/json' \
  -d '{"document_type":"POLICE_REPORT","file_name":"PR-2026-9912.pdf","file_url":"https://storage.googleapis.com/klemarklemer-claims-docs/PR-2026-9912.pdf"}' \
  | python3 -c '
import json,sys
c = json.load(sys.stdin)["data"]
r = c["recommendation"]
print("  stage      -> %s" % c["stage"])
print("  assigned   -> %s  (score %s)" % (c["current_officer"]["name"], c["assignment"]["total_score"]))
print("  recommends -> %s @ %s" % (r["outcome"], r["confidence"]))
print()
print("  " + r["reasons"])'
step

b "3 . The agent follows the record, not a script"
for loss in 90000 300 4200; do
  docker exec klemarklemer-postgres psql -U user -d klemarklemer_db \
    -c "update claims set estimated_loss=$loss where id=1;" >/dev/null 2>&1 || true
  printf '\n  estimated loss %s\n' "$loss"
  curl -s -X POST "$BASE/claim/1/assessment" | python3 -c '
import json,sys
r = json.load(sys.stdin)["data"]["recommendation"]
print("    -> %s @ %s" % (r["outcome"], r["confidence"]))
print("    %s" % r["reasons"][:220])'
done
step

b "4 . A human binds the decision - only this closes a claim"
curl -s -X POST "$BASE/claim/1/approval" -H 'Content-Type: application/json' \
  -d '{"action":"APPROVE","officer_id":1,"notes":"Verified against police report."}' \
  | python3 -c '
import json,sys
c = json.load(sys.stdin)["data"]
print("  stage=%s  status=%s  settlement=%s" % (c["stage"], c["status"], c["approved_amount"]))'
step

b "5 . The audit trail - every actor, every transition"
curl -s "$BASE/claim/1" | python3 -c '
import json,sys
for e in json.load(sys.stdin)["data"]["events"]:
    arrow = ""
    if e.get("new_stage"):
        arrow = "   %s -> %s" % (e.get("previous_stage") or "-", e["new_stage"])
    print("  %-28s %s%s" % (e["actor_name"], e["action"], arrow))'

b "Engine that produced the recommendation"
curl -s "$BASE/claim/1" | python3 -c '
import json,sys
for e in json.load(sys.stdin)["data"]["events"]:
    if e["action"] == "RECOMMENDATION_GENERATED":
        print("  " + json.loads(e["payload"]).get("engine",""))'
