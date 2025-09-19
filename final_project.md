# Final Project: W.H.Y

## Project Overview
The goal of this project is to build a **webhook-based authorization service** that integrates with **ArgoCD** to control deployment actions.  
The webhook will use **PlainID policies** to decide whether a requested deployment is authorized.  

This project demonstrates how **policy-based authorization** can be applied in CI/CD workflows, showcasing **DevOps practices** (GitHub, Kubernetes, ArgoCD, Helm, Docker) along with **PlainID’s authorization capabilities**.

---

## Architecture

### GitHub
- Stores source code and workflow definitions.  
- Triggers ArgoCD deployment pipeline.  

### K3s Cluster
- Local Kubernetes cluster hosting both ArgoCD and the custom webhook service.  

### ArgoCD
- Manages deployment of Helm charts into the cluster.  
- Configured with a webhook (pre-sync hook or external callout).  

### Custom Webhook Service
- Built using the **ArgoCD SDK**.  
- Packaged as a **Docker container** and deployed to the K3s cluster using Helm.  

**Responsibilities:**
- Receive deployment requests from ArgoCD.  
- Call **PlainID PDP** for an authorization decision.  
- If authorized → call ArgoCD API to proceed with the deployment.  
- If denied → block the deployment and log the event.  

### PlainID
- Acts as the **Policy Decision Point (PDP)**.  
- Evaluates *who can deploy, when, and under what conditions*.  
- Returns **ALLOW** or **DENY** back to the webhook.  

---

## Tech Stack
- **Kubernetes (K3s)** – Cluster runtime.  
- **ArgoCD** – GitOps-based deployment controller.  
- **GitHub** – Source control and CI triggers.  
- **PlainID PDP** – Policy decision engine.  
- **Webhook Service** –  
  - Language: Python (Flask/FastAPI), Node.js, or Go.  
  - Packaged with Docker.  
  - Deployed via Helm into K3s.  
- **ArgoCD SDK** – Programmatic interaction with ArgoCD.  

---

## Flow

```mermaid
sequenceDiagram
    participant Dev as Developer
    participant GH as GitHub
    participant Argo as ArgoCD
    participant WH as Webhook Service
    participant PDP as PlainID PDP

    Dev->>GH: Push code
    GH->>Argo: Trigger pipeline
    Argo->>WH: Send deployment request (pre-sync hook)
    WH->>PDP: Request authorization
    PDP-->>WH: ALLOW or DENY
    alt Authorized
        WH->>Argo: Proceed with deployment
        Argo->>K8s: Apply Helm changes
    else Denied
        WH->>Argo: Block deployment
        WH->>Dev: Log and notify
    end
