# W.H.Y Project

## Overview

The **W.H.Y Project** is a Kubernetes-based application that integrates ArgoCD for GitOps deployment, a pre-sync job to validate commit authorization, and a webhook service to interact with PlainID for access control.

This repository contains:

- Helm chart for deploying the application (`W.H.Y-chart`)
- Pre-sync job to fetch GitHub commit info and authorize ArgoCD syncs
- Webhook service to communicate with PlainID
- Dockerfiles and Go source code for the application and pre-sync job

---

## Prerequisites

Before deploying this project, make sure you have the following:

1. **Kubernetes Cluster** – any cluster where you have admin access.
2. **ArgoCD Installed** – ArgoCD must be installed and configured on the cluster.
3. **Helm v3+** – for managing the chart deployment.
4. **kubectl configured** – pointed to your cluster.
5. **GitHub Repository** – to store the project and fetch commits.
6. Optional: `jq` and `curl` installed in the environment if running the pre-sync job locally.

---

## Project Structure

.
├── README.md
├── W.H.Y-chart/
│ ├── Chart.yaml
│ ├── templates/
│ │ ├── argocd-application.yaml
│ │ ├── pre-sync-configmap.yaml
│ │ ├── pre-sync-job.yaml
│ │ ├── protectedapp-deployment.yaml
│ │ ├── protectedapp-service.yaml
│ │ ├── webhook-configmap.yaml
│ │ ├── webhook-deployment.yaml
│ │ ├── webhook-secrets.yaml
│ │ └── webhook-service.yaml
│ └── values.yaml
├── cmd/
│ └── main.go
├── config/
│ └── application.yaml
├── docker/
│ ├── W.H.Y-Dockerfile
│ └── pre-sync-Dockerfile
├── final_project.md
├── go.mod
├── go.sum
├── internal/
│ ├── config/
│ │ └── config.go
│ ├── plainid/
│ │ └── client.go
│ └── webhook/
│ └── server.go
└── k8s/
├── argo/
│ ├── alice_token
│ └── nana-myapp-argo-application.yaml
├── nana-myapp/
│ ├── deployment.yaml
│ ├── pre-sync-authorization-job.yaml
│ ├── pre-sync-configmap.yaml
│ └── service.yaml
└── webhook/
├── webhook-config.yaml
├── webhook-deployment.yaml
├── webhook-secrets.yaml
└── webhook-service.yaml


---

## Components

- **`protectedApp`**: Main application deployment.
- **`preSyncJob`**: Pre-sync job that fetches latest commit info and author from GitHub and calls the webhook for authorization.
- **`webhook`**: Service that interacts with PlainID and authorizes ArgoCD syncs.
- **`W.H.Y-chart`**: Helm chart with templates for all resources.
- **`Dockerfiles`**: For building `protectedApp` and `pre-sync` containers.
- **`ArgoCD Application`**: Template for syncing the Helm chart via ArgoCD.

---

## Helm Chart Deployment

### Values

Key parts of `values.yaml`:

```yaml
namespace: protectedapp

protectedApp:
  image: nanajanashia/argocd-app
  tag: "1.2"
  replicas: 2
  port: 8080

preSyncJob:
  configMapName: pre-sync-config
  image: yanivmendiuk/pre-sync-job
  tag: "1.0"
  env:
    appName: "protectedapp-argo-application"
    repoUrl: "https://github.com/YanivMendiuk/W.H.Y.git"
    webhookHost: "webhook.protectedapp.svc.cluster.local"
    webhookPort: 8181
  backoffLimit: 1
  terminationGracePeriodSeconds: 300

webhook:
  image: yanivmendiuk/webhook
  tag: "1.5"
  replicas: 1
  port: 8181
  config:
    entityTypeId: github
    resourceType: github_repos
    plainIdServiceName: my-plainid-paa-runtime
    plainIdNamespace: yaniv-namespace
    action: Access
  secrets:
    clientId: "<your-client-id>"
    clientSecret: "<your-client-secret>"

argocdApplication:
  enabled: true
  appName: protectedapp-argo-application
  argocdNamespace: argocd
  project: default
  repoURL: "https://github.com/YanivMendiuk/W.H.Y.git"
  path: "W.H.Y-chart"
  targetRevision: HEAD
  destinationServer: "https://kubernetes.default.svc"
  destinationNamespace: "{{ .Values.namespace }}"
  syncOptions:
    - CreateNamespace=true
  automated:
    selfHeal: true
    prune: true
```
---

## Deployment Overview

```bash
helm upgrade --install protectedapp ./W.H.Y-chart \
  --namespace protectedapp \
  --create-namespace \
  --set webhook.secrets.clientId="YOUR_CLIENT_ID" \
  --set webhook.secrets.clientSecret="YOUR_CLIENT_SECRET"
```

This will deploy:

- **protectedApp** Deployment and Service
- **preSyncJob** to fetch commit info and authorize
- **webhook** Deployment and Service with secrets

---

## ArgoCD Integration

If `argocdApplication.enabled` is set to `true`, the Helm chart will also deploy an **ArgoCD Application**:

- **Path**: `W.H.Y-chart`
- **Target revision**: `HEAD`
- **Destination namespace**: `protectedapp`
- **Sync policy**: automated with self-heal and pruning
- Optional Helm `valueFiles` or `parameters` can be configured

> **Note:** ArgoCD does not pass local `--set` secrets automatically. Ensure secrets are injected either via Helm or a Kubernetes Secret for production.

---

## Pre-sync Job Details

The **pre-sync job** is configured as a Kubernetes Job with an ArgoCD `PreSync` hook:

- Fetches the latest commit SHA and author email from GitHub
- Sends authorization request to the webhook
- Job fails if the webhook returns non-200 response

This ensures that only authorized commits are synced by ArgoCD.

---

## Webhook Service

- Reads configuration from a ConfigMap
- Uses Kubernetes Secret for `clientId` and `clientSecret`
- Communicates with PlainID runtime API for authorization
- Injected as environment variables into the container

---

## Best Practices

- Keep secrets out of GitHub and use Helm `--set` or Kubernetes Secrets
- Always use the same namespace for all components to avoid misconfigurations
- If you modify the chart locally and want to force sync with ArgoCD:

```bash
argocd app sync protectedapp-argo-application --force
```
- Test locally using Helm before pushing changes.
- Document any new parameters in values.yaml.



---

This README includes:  

- Overview of the project  
- Prerequisites (ArgoCD, Helm, kubectl, etc.)  
- Clear explanation of all components (`protectedApp`, `preSyncJob`, `webhook`, ArgoCD app)  
- Detailed instructions for Helm deployment with secrets  
- ArgoCD sync details  
- Pre-sync and webhook explanation  
- Useful commands for working with the cluster  
