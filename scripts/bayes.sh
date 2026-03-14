#!/usr/bin/env bash
set -euo pipefail

if [ $# -lt 1 ]; then
  echo "Usage: $0 <json_file>" >&2
  exit 1
fi

file="$1"
scale=10

initial_prior=$(jq -r '.initial_prior_a' "$file")
desc_a=$(jq -r '.desc_a // empty' "$file")
count=$(jq '.iteration_history | length' "$file")

if [ "$count" -eq 0 ]; then
  echo "No iteration history."
  exit 0
fi

# 标题
if [ -n "$desc_a" ] && [ "$desc_a" != "A" ]; then
  echo "Bayesian Iteration History  (A = \"$desc_a\")"
else
  echo "Bayesian Iteration History"
fi
echo "============================================================"
echo

prior="$initial_prior"

for i in $(seq 0 $((count - 1))); do
  record=$(jq ".iteration_history[$i]" "$file")
  la=$(echo "$record" | jq -r '.likelihood_a')
  lna=$(echo "$record" | jq -r '.likelihood_not_a')
  desc_b=$(echo "$record" | jq -r '.desc_b // empty')

  # P(B) = P(B|A)*P(A) + P(B|¬A)*(1-P(A))
  pb=$(echo "scale=$scale; $la * $prior + $lna * (1 - $prior)" | bc)
  # P(A|B) = P(B|A)*P(A) / P(B)
  posterior=$(echo "scale=$scale; $la * $prior / $pb" | bc)

  # 转百分比显示
  prior_pct=$(echo "scale=4; $prior * 100" | bc)
  la_pct=$(echo "scale=4; $la * 100" | bc)
  lna_pct=$(echo "scale=4; $lna * 100" | bc)
  pb_pct=$(echo "scale=4; $pb * 100" | bc)
  post_pct=$(echo "scale=4; $posterior * 100" | bc)

  # 迭代编号
  num=$((i + 1))
  header="#$num"
  if [ -n "$desc_b" ] && [ "$desc_b" != "B" ]; then
    header="$header  B = \"$desc_b\""
  fi
  echo "$header"
  echo "----------------------------------------"
  echo "  P(A)   = ${prior_pct}%"
  echo "  P(B|A) = ${la_pct}%,  P(B|¬A) = ${lna_pct}%"
  echo "  P(B)   = P(B|A)·P(A) + P(B|¬A)·P(¬A) = ${pb_pct}%"
  echo "  P(A|B) = P(B|A)·P(A) / P(B) = ${post_pct}%"
  echo

  # 下一轮的先验 = 本轮后验
  prior="$posterior"
done

echo "============================================================"
if [ -n "$desc_a" ] && [ "$desc_a" != "A" ]; then
  echo "After $count iteration(s), the probability of \"$desc_a\" is ${post_pct}%."
else
  echo "After $count iteration(s), P(A) = ${post_pct}%."
fi
