#!/bin/bash
set -euo pipefail

KIND_CLUSTER=${KIND_CLUSTER:-llama}

kind create cluster --name "$KIND_CLUSTER"

# build image
IMAGE=llama-operator:latest
docker build -t $IMAGE ..
kind load docker-image --name "$KIND_CLUSTER" $IMAGE

kubectl apply -f ../llama-operator/config/crd.yaml
kubectl apply -f ../llama-operator/config/rbac.yaml
kubectl apply -f ../llama-operator/config/operator.yaml
kubectl apply -f ../llama-operator/config/sample.yaml

kubectl rollout status deployment/llama-operator -n default
