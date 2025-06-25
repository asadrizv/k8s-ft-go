#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"

kubectl apply -f "$REPO_ROOT/config/crd/bases"
kubectl apply -f "$REPO_ROOT/config/rbac"
kubectl apply -f "$REPO_ROOT/config/manager"
