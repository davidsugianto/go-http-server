# Kubernetes Learning & Teaching Module — VirtualBox Edition
### A Hands-On Curriculum for Beginner & Intermediate Learners — Everything Runs Inside VirtualBox VMs, Built Entirely on Free Tools

---

## 1. How to Use This Module

**For self-learners:** Work through the tracks in order. Each module has learning objectives, core topics, a **detailed step-by-step lab**, and a "you're done when" checkpoint. Budget 4–6 hours per module.

**For instructors:** Each module maps to one week (2 hrs lecture + 2–4 hrs lab). The VirtualBox approach means every student has an identical, disposable environment — distribute a pre-built `.ova` appliance export to skip setup time, and use snapshots so a broken lab is a 10-second rollback, never a lost session.

**Why VirtualBox for everything?**
- Identical environment for every learner regardless of host OS (Windows/Linux/Intel Mac).
- **Snapshots** = instant undo for any mistake. This changes how fearlessly students experiment.
- The intermediate track builds a *real multi-node cluster* across 3 VMs with kubeadm — the closest free experience to production.
- Nothing installed on the host machine except VirtualBox itself.

**Total duration:**
- Beginner Track: ~6 weeks (Part 0 setup + Modules B1–B6 + Capstone 1) — runs in **one VM** with minikube.
- Intermediate Track: ~8 weeks (Modules I1–I7 + Capstone 2) — runs on a **3-VM kubeadm cluster**.

---

## 2. Lab Architecture

```
┌─────────────────────────────  YOUR HOST MACHINE  ─────────────────────────────┐
│  Browser ──► http://192.168.56.x        Terminal ──► ssh k8s@192.168.56.x     │
│                                                                               │
│  ┌──────────────────────────  VIRTUALBOX  ───────────────────────────────┐    │
│  │                                                                       │    │
│  │  BEGINNER TRACK (one VM)          INTERMEDIATE TRACK (three VMs)      │    │
│  │  ┌──────────────────────┐         ┌───────────┐ ┌────────┐ ┌────────┐ │    │
│  │  │ k8s-lab              │         │ k8s-cp    │ │ k8s-w1 │ │ k8s-w2 │ │    │
│  │  │ 192.168.56.10        │         │ .56.10    │ │ .56.11 │ │ .56.12 │ │    │
│  │  │ Ubuntu + Docker      │         │ control   │ │ worker │ │ worker │ │    │
│  │  │ + minikube cluster   │         │ plane     │ │        │ │        │ │    │
│  │  └──────────────────────┘         └───────────┘ └────────┘ └────────┘ │    │
│  │                                     built with kubeadm + Calico       │    │
│  │  Adapter 1: NAT (internet access)                                     │    │
│  │  Adapter 2: Host-Only 192.168.56.0/24 (host ↔ VM communication)       │    │
│  └───────────────────────────────────────────────────────────────────────┘    │
└───────────────────────────────────────────────────────────────────────────────┘
```

**Host machine requirements**

| Track | Host RAM | Host CPU | Free disk |
|---|---|---|---|
| Beginner (1 VM) | 8 GB min, 16 GB comfortable | 4 cores | 60 GB |
| Intermediate (3 VMs) | 16 GB recommended | 4–8 cores | 120 GB |
| Fallback if host has only 8 GB | Run intermediate labs on the single beginner VM using a multi-node **kind** cluster inside it | | |

> ⚠️ **Apple Silicon (M1–M4) note:** classic VirtualBox targets x86. Use VirtualBox 7.1+ (Arm host support, with **Ubuntu Server for ARM** ISO), or substitute the identical VM steps in **UTM** or **Multipass** (both free). Every command *inside* the VMs is unchanged.

---

## 3. Part 0 — Build the Lab Environment (Detailed Step-by-Step)

### 0.1 Install VirtualBox on the host

1. Download VirtualBox from **https://www.virtualbox.org/wiki/Downloads** — pick your host platform (Windows / Linux / macOS).
2. Run the installer with defaults. On Windows, approve the network-driver prompts (needed for host-only networking).
3. Optionally install the **Extension Pack** from the same page (free for personal/educational use).
4. Verify: launch VirtualBox — you should see the empty VirtualBox Manager window.

> **Windows users:** if VMs later crash or run extremely slowly, disable Hyper-V conflicts: `Turn Windows features on or off` → untick *Hyper-V*, *Virtual Machine Platform*, *Windows Hypervisor Platform* → reboot. (Skip if you also need WSL2/Docker Desktop; VirtualBox 7 can coexist, just slower.)

### 0.2 Download the Ubuntu Server ISO

1. Go to **https://ubuntu.com/download/server** and download **Ubuntu Server 24.04 LTS** (~2.6 GB).
2. Server, not Desktop: no GUI means the VM needs far less RAM, and you'll work over SSH like a real admin.

### 0.3 Create the Host-Only network (do this once)

This gives your VMs fixed IPs your host browser and terminal can reach.

1. VirtualBox Manager → **File → Tools → Network Manager** (older versions: File → Host Network Manager).
2. **Host-only Networks** tab → **Create**. You'll get `vboxnet0` (Linux/Mac) or `VirtualBox Host-Only Ethernet Adapter` (Windows).
3. Select it → **Properties**:
   - IPv4 Address: `192.168.56.1`, Mask: `255.255.255.0` (usually the default).
   - **DHCP Server tab: untick "Enable Server"** — we'll use static IPs so cluster addresses never change.
4. Apply.

### 0.4 Create the base VM

1. VirtualBox Manager → **New**:
   - **Name:** `k8s-lab` · **Type:** Linux · **Version:** Ubuntu (64-bit)
   - **ISO Image:** select the Ubuntu Server ISO. Tick **"Skip Unattended Installation"** (we want the real installer — it's a learning experience).
2. **Hardware:**
   - Memory: **8192 MB** (minimum 4096)
   - Processors: **4** (minimum 2)
   - Tick **Enable EFI** (optional, either works)
3. **Hard Disk:** Create a virtual hard disk, **VDI**, dynamically allocated, **50 GB**.
4. Click **Finish**, then open **Settings** for the new VM before booting:
   - **System → Processor:** tick *Enable PAE/NX*.
   - **Network → Adapter 1:** Attached to **NAT** (this is the VM's internet access).
   - **Network → Adapter 2:** tick *Enable Network Adapter* → Attached to **Host-only Adapter** → select `vboxnet0`.
   - **Audio / USB:** untick (not needed, saves resources).

### 0.5 Install Ubuntu inside the VM

1. **Start** the VM; choose *Try or Install Ubuntu Server*.
2. Walk the installer:
   - Language / keyboard: your choice.
   - **Network:** you'll see two interfaces — `enp0s3` (NAT, gets 10.0.2.15 via DHCP — leave it) and `enp0s8` (host-only, "not connected" or no address — leave it for now; we configure it after install). Continue.
   - Proxy: blank. Mirror: default.
   - **Storage:** *Use entire disk* (it's virtual — nothing on your host is touched). Confirm.
   - **Profile:** name anything; **server name:** `k8s-lab`; **username:** `k8s`; pick a password you'll remember.
   - **SSH:** ✅ **tick "Install OpenSSH server"** — critical, this is how you'll work.
   - Featured snaps: select **none** (we install Docker properly ourselves).
3. Wait for install → **Reboot Now**. If it complains about the CD, VirtualBox usually auto-ejects; press Enter.
4. Log in at the VM console as `k8s`.

### 0.6 Configure the static host-only IP

At the VM console:

```bash
# Confirm interface names (expect enp0s3 = NAT with 10.0.2.15, enp0s8 = no IP yet)
ip -br a

# Create a netplan file for the host-only adapter
sudo tee /etc/netplan/60-hostonly.yaml > /dev/null <<'EOF'
network:
  version: 2
  ethernets:
    enp0s8:
      dhcp4: no
      addresses: [192.168.56.10/24]
EOF

sudo chmod 600 /etc/netplan/60-hostonly.yaml
sudo netplan apply
ip -br a    # enp0s8 should now show 192.168.56.10/24
```

**From your HOST machine**, verify and switch to a comfortable terminal:

```bash
ping 192.168.56.10        # should reply
ssh k8s@192.168.56.10     # log in — work from here from now on
```

> Windows: use the built-in `ssh` in PowerShell/Terminal, or PuTTY. You now get copy-paste, scrollback, and multiple tabs — never suffer the VM console again.

### 0.7 Install the tool chain inside the VM

All commands run **inside the VM over SSH**.

**System update:**
```bash
sudo apt-get update && sudo apt-get upgrade -y
```

**Docker Engine (official repository):**
```bash
sudo apt-get install -y ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] \
https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo $VERSION_CODENAME) stable" \
| sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Let your user run docker without sudo
sudo usermod -aG docker $USER
exit        # log out, then ssh back in so the group applies
```
```bash
ssh k8s@192.168.56.10
docker run hello-world     # must succeed without sudo
```

**kubectl:**
```bash
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
kubectl version --client
```

**minikube:**
```bash
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube && rm minikube-linux-amd64
```

**Helm, k9s, and quality-of-life:**
```bash
# Helm
curl -fsSL https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

# k9s terminal UI
curl -LO https://github.com/derailed/k9s/releases/latest/download/k9s_linux_amd64.deb
sudo dpkg -i k9s_linux_amd64.deb && rm k9s_linux_amd64.deb

# kubectl alias + autocompletion (huge time-saver)
cat >> ~/.bashrc <<'EOF'
alias k=kubectl
source <(kubectl completion bash)
complete -o default -F __start_kubectl k
EOF
source ~/.bashrc
```

### 0.8 The snapshot workflow (your superpower)

Shut the VM down (`sudo poweroff`), then in VirtualBox Manager select `k8s-lab` → **Snapshots** → **Take**:

- Name it **`00-clean-tools`** — "Ubuntu + Docker + kubectl + minikube + helm, no cluster yet".

**Rules of thumb from here on:**
- Take a snapshot **before every module** (`B2-start`, `B3-start`, …).
- Broke something beyond repair? **Restore** the snapshot — 10 seconds, fresh start.
- Snapshots are also how instructors "reset the room" between cohorts.
- CLI equivalent: `VBoxManage snapshot k8s-lab take "B2-start"`.

### 0.9 Start your first cluster

Boot the VM, SSH in:

```bash
minikube start --driver=docker --cpus=2 --memory=4096
# If your VM only has 4 GB RAM, use: --memory=2200

kubectl get nodes
kubectl cluster-info
minikube addons enable metrics-server
minikube addons enable dashboard
```

Expected: one node named `minikube`, status `Ready`. **Take a snapshot: `01-cluster-running`.**

### 0.10 Reaching your apps from the host browser (read this — it matters all course)

Because minikube runs *inside* the VM (as a Docker container), its NodePorts are only visible inside the VM. Two clean patterns:

**Pattern A — kubectl port-forward bound to all interfaces (use this 90% of the time):**
```bash
# inside the VM: forward local port 8080 → service port 80, listening on ALL VM interfaces
kubectl port-forward --address 0.0.0.0 service/<svc-name> 8080:80
```
Then on your **host browser**: `http://192.168.56.10:8080` ✅

**Pattern B — SSH tunnel from the host (no cluster changes at all):**
```bash
# on the HOST
ssh -L 8080:localhost:8080 k8s@192.168.56.10
# then inside that session: kubectl port-forward service/<svc> 8080:80
# host browser → http://localhost:8080
```

> In the **intermediate track** this limitation disappears: the kubeadm nodes' IPs (192.168.56.10–12) are directly reachable, so NodePorts and MetalLB LoadBalancer IPs work straight from your host browser.

---

## 4. The Free Tool Stack (VirtualBox Edition)

| Layer | Tool | Cost |
|---|---|---|
| Virtualization | **VirtualBox 7.x** | Free (GPL); Extension Pack free for personal/education |
| Guest OS | **Ubuntu Server 24.04 LTS** | Free |
| Container runtime | **Docker Engine** (beginner VM), **containerd** (kubeadm nodes) | Free/open source |
| Beginner cluster | **minikube** (docker driver, inside VM) | Open source |
| Intermediate cluster | **kubeadm** — real 3-node cluster across VMs | Built into Kubernetes |
| CNI / NetworkPolicy | **Calico** | Open source |
| LoadBalancer on VMs | **MetalLB** | Open source |
| CLI & UX | **kubectl, k9s, helm, kubectx/kubens** | Open source |
| Ingress & TLS | **NGINX Ingress Controller, cert-manager** | Open source |
| Observability | **kube-prometheus-stack (Prometheus+Grafana), Loki** | Open source |
| CI/CD & GitOps | **GitHub Actions (free tier), Argo CD, ghcr.io / Docker Hub** | Free tiers |
| Security | **Trivy, kube-bench**, native RBAC & NetworkPolicy | Open source |
| Load testing | **k6** or **hey** | Open source |
| Zero-install fallback | **Killercoda, Play with Kubernetes** | Free |

---
## 5. Beginner Track (Weeks 1–6) — Detailed Labs
> All labs run inside the **k8s-lab VM** over SSH. Snapshot before each module. `k` is your alias for `kubectl`.

### Module B1 — Why Kubernetes? Architecture & First Contact
**Objectives:** explain what K8s solves; name every control-plane and node component; run a first workload.
**Topics:** orchestration vs plain containers · API server, etcd, scheduler, controller-manager · kubelet, kube-proxy, runtime · declarative model & the reconciliation loop.

**Lab B1 — step by step:**
```bash
# 1. Explore the cluster you built in Part 0
kubectl get nodes -o wide
kubectl describe node minikube        # find capacity, addresses, system pods
kubectl get pods -n kube-system       # meet the control plane itself

# 2. Run your first pod
kubectl run hello --image=nginx
kubectl get pods -w                   # watch it go Pending → ContainerCreating → Running (Ctrl-C to stop)
kubectl describe pod hello            # read the Events section bottom-up — this habit will save you weekly

# 3. Look inside
kubectl logs hello
kubectl exec -it hello -- bash        # you are now INSIDE the container
ls /usr/share/nginx/html && exit

# 4. See it in a browser (Pattern A from §0.10)
kubectl port-forward --address 0.0.0.0 pod/hello 8080:80
# Host browser → http://192.168.56.10:8080  → nginx welcome page. Ctrl-C when done.

# 5. The punchline: pods alone don't heal
kubectl delete pod hello
kubectl get pods                      # gone forever — nothing recreated it. Why? → Module B2
```
**Optional — the dashboard from your host:**
```bash
kubectl proxy --address 0.0.0.0 --accept-hosts='.*' &
# Host browser → http://192.168.56.10:8001/api/v1/namespaces/kubernetes-dashboard/services/http:kubernetes-dashboard:/proxy/
```
**✅ Done when:** you can sketch the architecture from memory and narrate exactly what happened between `kubectl run` and `Running`.

---

### Module B2 — Pods, ReplicaSets & Deployments
**Objectives:** author manifests from scratch; roll out, update, and roll back with Deployments.
**Topics:** pod spec anatomy · labels & selectors · ReplicaSets · rolling updates & rollbacks · `apply` vs `create` · `--dry-run=client -o yaml` as a manifest generator.

**Lab B2 — step by step:**
```bash
mkdir -p ~/labs/b2 && cd ~/labs/b2
```
1. **Write your first Deployment by hand** (`vim deploy.yaml` — typing it builds recall):
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web
  labels: {app: web}
spec:
  replicas: 3
  selector:
    matchLabels: {app: web}
  template:
    metadata:
      labels: {app: web}
    spec:
      containers:
      - name: nginx
        image: nginx:1.27
        ports:
        - containerPort: 80
```
```bash
kubectl apply -f deploy.yaml
kubectl get deploy,rs,pods -o wide     # see the three-layer ownership: Deployment → ReplicaSet → Pods
```
2. **Prove self-healing** — kill a pod, watch its replacement appear in seconds:
```bash
kubectl delete pod $(kubectl get pods -l app=web -o name | head -1)
kubectl get pods -l app=web -w
```
3. **Rolling update** — bump the image and watch the choreography:
```bash
kubectl set image deployment/web nginx=nginx:1.28
kubectl rollout status deployment/web
kubectl get rs                         # old RS scaled to 0, new RS to 3
kubectl rollout history deployment/web
```
4. **Break it, then time-travel back:**
```bash
kubectl set image deployment/web nginx=nginx:definitely-not-a-tag
kubectl get pods                       # ImagePullBackOff — diagnose with: kubectl describe pod <name>
kubectl rollout undo deployment/web
kubectl rollout status deployment/web  # healthy again
```
5. **Generator trick** you'll use forever:
```bash
kubectl create deployment demo --image=nginx --replicas=2 --dry-run=client -o yaml > generated.yaml
```
**✅ Done when:** blank file → running → updated → broken → rolled back, without looking anything up. **Snapshot: `B2-done`.**

---

### Module B3 — Services & Networking Basics
**Objectives:** expose apps inside and outside the cluster; understand service discovery.
**Topics:** ephemeral pod IPs · ClusterIP / NodePort / LoadBalancer · CoreDNS names · endpoints · port vs targetPort vs nodePort.

**Lab B3 — step by step:**
1. **An app that identifies itself** (each replica reports its own hostname):
```yaml
# ~/labs/b3/echo.yaml
apiVersion: apps/v1
kind: Deployment
metadata: {name: echo}
spec:
  replicas: 3
  selector: {matchLabels: {app: echo}}
  template:
    metadata: {labels: {app: echo}}
    spec:
      containers:
      - name: echo
        image: registry.k8s.io/e2e-test-images/agnhost:2.47
        args: ["netexec", "--http-port=8080"]
        ports: [{containerPort: 8080}]
---
apiVersion: v1
kind: Service
metadata: {name: echo-svc}
spec:
  selector: {app: echo}
  ports: [{port: 80, targetPort: 8080}]
```
```bash
kubectl apply -f ~/labs/b3/echo.yaml
kubectl get svc echo-svc               # note the stable ClusterIP
kubectl get endpoints echo-svc         # the 3 pod IPs behind it
```
2. **Consume it like another microservice would** — via DNS, from a throwaway pod:
```bash
kubectl run tmp --rm -it --image=busybox:1.36 -- sh
# inside:
wget -qO- http://echo-svc/hostname     # run 5–6 times → different pod names = load balancing
wget -qO- http://echo-svc.default.svc.cluster.local/hostname   # the full DNS form
exit
```
3. **NodePort** — expose outside the cluster:
```bash
kubectl patch svc echo-svc -p '{"spec":{"type":"NodePort"}}'
kubectl get svc echo-svc               # note the 3xxxx port
curl $(minikube ip):<nodePort>/hostname     # works INSIDE the VM
# For the host browser, use the standard pattern:
kubectl port-forward --address 0.0.0.0 svc/echo-svc 8080:80
# host → http://192.168.56.10:8080/hostname — refresh to see balancing
```
4. **Discussion checkpoint:** trace a request host-browser → VM:8080 → port-forward → Service → pod, naming each hop.
**✅ Done when:** you can explain every hop and why pod IPs must never be hard-coded.

---

### Module B4 — Configuration: ConfigMaps, Secrets & Environment
**Objectives:** separate config from code; inject config as env vars and files; handle credentials.
**Topics:** ConfigMaps (env + volume) · Secrets & why base64 ≠ encryption · `envFrom` · redeploying on config change.

**Lab B4 — step by step:**
```bash
mkdir -p ~/labs/b4 && cd ~/labs/b4
# 1. Create config two ways
kubectl create configmap app-config --from-literal=APP_COLOR=blue --from-literal=APP_MODE=dev
kubectl create secret generic db-secret --from-literal=DB_PASSWORD='S3cr3t!'
```
2. **Consume both in a visible way** — agnhost echoes its env via HTTP:
```yaml
# color-app.yaml
apiVersion: apps/v1
kind: Deployment
metadata: {name: color-app}
spec:
  replicas: 1
  selector: {matchLabels: {app: color}}
  template:
    metadata: {labels: {app: color}}
    spec:
      containers:
      - name: app
        image: registry.k8s.io/e2e-test-images/agnhost:2.47
        args: ["netexec", "--http-port=8080"]
        envFrom:
        - configMapRef: {name: app-config}
        env:
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef: {name: db-secret, key: DB_PASSWORD}
```
```bash
kubectl apply -f color-app.yaml
kubectl exec deploy/color-app -- env | grep -E 'APP_|DB_'
```
3. **Change config → observe that running pods DON'T see it** (env is set at start):
```bash
kubectl create configmap app-config --from-literal=APP_COLOR=red --from-literal=APP_MODE=dev -o yaml --dry-run=client | kubectl apply -f -
kubectl exec deploy/color-app -- env | grep APP_COLOR    # still blue!
kubectl rollout restart deployment/color-app             # the standard fix
kubectl exec deploy/color-app -- env | grep APP_COLOR    # now red
```
4. **ConfigMap as a mounted file** — serve a custom page with stock nginx:
```bash
kubectl create configmap web-content --from-literal=index.html='<h1>Config-driven page v1</h1>'
```
```yaml
# web.yaml — key part of the pod template:
      containers:
      - name: nginx
        image: nginx:1.27
        volumeMounts:
        - {name: content, mountPath: /usr/share/nginx/html}
      volumes:
      - name: content
        configMap: {name: web-content}
```
```bash
kubectl apply -f web.yaml
kubectl port-forward --address 0.0.0.0 deploy/web 8080:80   # host browser shows your page
```
5. **The secret "reveal"** — and the lesson:
```bash
kubectl get secret db-secret -o jsonpath='{.data.DB_PASSWORD}' | base64 -d; echo
# Discussion: base64 is encoding, not encryption. Real remedies: RBAC on secrets,
# encryption-at-rest, external managers (sealed-secrets/SOPS — see Intermediate stretch goals).
```
**✅ Done when:** the same image serves different content/behavior in two namespaces purely via config (`kubectl create ns staging` and repeat with different values).

---

### Module B5 — Storage: Volumes, PVs & PVCs
**Objectives:** persist data beyond pod life; trace claim → volume → class.
**Topics:** emptyDir vs PVC · PersistentVolume, PersistentVolumeClaim, StorageClass · dynamic provisioning · access modes & reclaim policy.

**Lab B5 — step by step:**
1. **Postgres with NO volume — watch data die:**
```yaml
# ~/labs/b5/pg-novol.yaml
apiVersion: apps/v1
kind: Deployment
metadata: {name: pg}
spec:
  replicas: 1
  selector: {matchLabels: {app: pg}}
  template:
    metadata: {labels: {app: pg}}
    spec:
      containers:
      - name: postgres
        image: postgres:16
        env:
        - name: POSTGRES_PASSWORD
          valueFrom: {secretKeyRef: {name: db-secret, key: DB_PASSWORD}}
```
```bash
kubectl apply -f pg-novol.yaml && kubectl wait --for=condition=ready pod -l app=pg
kubectl exec -it deploy/pg -- psql -U postgres -c "CREATE TABLE t(x int); INSERT INTO t VALUES (42);"
kubectl delete pod -l app=pg           # simulate a crash
kubectl wait --for=condition=ready pod -l app=pg
kubectl exec -it deploy/pg -- psql -U postgres -c "SELECT * FROM t;"   # ERROR: table gone
```
2. **Add a PVC** (minikube's default StorageClass provisions automatically):
```yaml
# pvc.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata: {name: pg-data}
spec:
  accessModes: [ReadWriteOnce]
  resources: {requests: {storage: 1Gi}}
```
Add to the pod template:
```yaml
        volumeMounts:
        - {name: data, mountPath: /var/lib/postgresql/data}
        env:
        - {name: PGDATA, value: /var/lib/postgresql/data/pgdata}
      volumes:
      - name: data
        persistentVolumeClaim: {claimName: pg-data}
```
```bash
kubectl apply -f pvc.yaml -f pg-vol.yaml
# repeat the insert → delete pod → select test: the row SURVIVES ✅
```
3. **Trace the chain:**
```bash
kubectl get pvc,pv,sc
kubectl describe pvc pg-data           # who bound it, to what, provisioned by which class
```
**✅ Done when:** you can explain to a rubber duck who creates the PV, who creates the PVC, and what the StorageClass did — and demo data surviving a pod kill.

---

### Module B6 — Operating & Debugging: Namespaces, Probes, Limits & the Toolbox
**Objectives:** structure clusters; make workloads production-polite; develop a debugging reflex.
**Topics:** namespaces · liveness/readiness/startup probes · requests & limits, QoS · `logs / describe / exec / events / port-forward / get -o yaml`.

**Lab B6 — step by step:**
1. **Production-polite pod** — add to your B4 app's container spec:
```yaml
        readinessProbe:
          httpGet: {path: /healthz, port: 8080}
          initialDelaySeconds: 3
          periodSeconds: 5
        livenessProbe:
          httpGet: {path: /healthz, port: 8080}
          initialDelaySeconds: 10
          periodSeconds: 10
        resources:
          requests: {cpu: 100m, memory: 64Mi}
          limits: {cpu: 250m, memory: 128Mi}
```
Apply, then `kubectl describe pod` → see probe config; `kubectl get events --sort-by=.lastTimestamp` after any hiccup.
2. **The Broken Five** — deploy these and fix each using only kubectl (instructors: hand out as files; self-learners: type them, break them on purpose):
   1. Image `nginx:1.27-typo` → *ImagePullBackOff* → fix tag.
   2. Liveness probe on the wrong port → *CrashLoop of restarts* → read `describe`, fix port.
   3. `requests: {memory: 64Gi}` → *Pending / Insufficient memory* → read `kubectl describe pod` scheduler events.
   4. Service selector `app: webb` (typo) → app "up" but curl fails → `kubectl get endpoints` shows none → fix selector.
   5. `envFrom` a ConfigMap that doesn't exist → *CreateContainerConfigError* → create it or fix the name.
3. **Namespaces & context hygiene:**
```bash
kubectl create ns dev; kubectl create ns staging
kubectl apply -n dev -f color-app.yaml
kubectl apply -n staging -f color-app.yaml   # different ConfigMap values per ns
kubectl config set-context --current --namespace=dev
```
4. Spend 20 minutes in **k9s** (`k9s` → `:pods`, `l` logs, `d` describe, `s` shell) — it will halve your debugging time all course.
**✅ Done when:** all five broken deployments fixed in ≤30 minutes, narrating your reasoning.

---

### 🏆 Capstone Project 1 — Multi-Tier App on the VM Cluster
**App:** Docker's **Example Voting App** (github.com/dockersamples/example-voting-app) — vote frontend (Python) → Redis → worker (.NET) → Postgres → results frontend (Node).

**Step-by-step skeleton:**
```bash
git clone https://github.com/dockersamples/example-voting-app ~/labs/capstone1
cd ~/labs/capstone1/k8s-specifications   # reference manifests — you will REWRITE, not copy
kubectl create ns vote
```
Requirements checklist (build your own manifests in a Git repo):
1. Deployments for vote, redis, worker, result, db — db with a **PVC** (B5).
2. ClusterIP services for redis & db; the two frontends reachable from the **host browser** via `port-forward --address 0.0.0.0` (vote → 8080, result → 8081).
3. All tunables in **ConfigMaps**, credentials in **Secrets** — nothing hard-coded (B4).
4. Probes + requests/limits on every container (B6).
5. Everything in the `vote` namespace; repo has README + architecture diagram.
6. **Live demo:** vote from the host browser → kill the redis pod and a frontend pod mid-demo → app self-heals, tallies survive (db PVC).

**Stretch:** HPA on the vote frontend (needs metrics-server, already enabled) · a `kubectl`-only runbook for the Broken-Five failure classes · export your VM as an `.ova` and hand your whole working environment to a friend.

**Snapshot when done: `capstone1-complete` — the Intermediate track builds new VMs, but keep this one as your sandbox.**

---
## 6. Intermediate Setup — Build a Real 3-Node Cluster in VirtualBox (kubeadm)
> This is itself the first intermediate lab: you assemble Kubernetes the way cloud providers do. Budget one full session.

### 6.1 Create the three VMs

Fastest route — **clone** your clean beginner base instead of reinstalling Ubuntu:

1. Power off `k8s-lab`, restore snapshot **`00-clean-tools`** (pre-cluster, pre-minikube state).
2. Right-click → **Clone** → name `k8s-cp` → **Full clone** → MAC Address Policy: **Generate new MAC addresses for all adapters** (critical — duplicate MACs break networking). Repeat for `k8s-w1` and `k8s-w2`.
3. Adjust resources per VM (Settings → System):

| VM | vCPU | RAM | Role | Host-only IP |
|---|---|---|---|---|
| k8s-cp | 2 | 4096 MB | control plane | 192.168.56.10 |
| k8s-w1 | 2 | 3072 MB | worker | 192.168.56.11 |
| k8s-w2 | 2 | 3072 MB | worker | 192.168.56.12 |

> Only 8 GB host RAM? Run cp + one worker (drop k8s-w2), or do the whole intermediate track on the beginner VM with `kind create cluster --config` (3-node config file) — every Kubernetes command below is identical.

4. Boot each clone **one at a time** and fix its identity (each was cloned from the same machine):
```bash
# ON k8s-cp (repeat with the right name/IP on w1 and w2)
sudo hostnamectl set-hostname k8s-cp
sudo sed -i 's/k8s-lab/k8s-cp/g' /etc/hosts
sudo sed -i 's/192.168.56.10/192.168.56.10/' /etc/netplan/60-hostonly.yaml   # w1 → .11, w2 → .12
sudo netplan apply
sudo reboot
```
5. Give every node a map of the cluster — append to `/etc/hosts` on **all three**:
```bash
cat <<'EOF' | sudo tee -a /etc/hosts
192.168.56.10 k8s-cp
192.168.56.11 k8s-w1
192.168.56.12 k8s-w2
EOF
```
6. From the host, open **three terminal tabs**: `ssh k8s@192.168.56.10 / .11 / .12`.

### 6.2 Prepare ALL three nodes (run everywhere)

```bash
# 1. Swap off — kubelet refuses to run with swap
sudo swapoff -a
sudo sed -i '/\sswap\s/ s/^/#/' /etc/fstab
sudo sed -i '/swap.img/ s/^/#/' /etc/fstab

# 2. Kernel modules & sysctls for container networking
cat <<EOF | sudo tee /etc/modules-load.d/k8s.conf
overlay
br_netfilter
EOF
sudo modprobe overlay && sudo modprobe br_netfilter

cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-iptables  = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.ipv4.ip_forward                 = 1
EOF
sudo sysctl --system

# 3. containerd as the runtime (already installed from the docker repo in the base image)
sudo mkdir -p /etc/containerd
containerd config default | sudo tee /etc/containerd/config.toml > /dev/null
sudo sed -i 's/SystemdCgroup = false/SystemdCgroup = true/' /etc/containerd/config.toml
sudo systemctl restart containerd && sudo systemctl enable containerd

# 4. kubeadm, kubelet, kubectl from the official pkgs.k8s.io repo
#    (v1.33 shown — swap in the current stable minor from kubernetes.io/releases)
sudo apt-get install -y apt-transport-https ca-certificates curl gpg
curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.33/deb/Release.key \
  | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.33/deb/ /' \
  | sudo tee /etc/apt/sources.list.d/kubernetes.list
sudo apt-get update
sudo apt-get install -y kubelet kubeadm kubectl
sudo apt-mark hold kubelet kubeadm kubectl

# 5. THE classic VirtualBox gotcha: two NICs → kubelet may advertise the NAT IP (10.0.2.15)
#    Pin each node to its host-only IP (use .10 / .11 / .12 per node!):
echo "KUBELET_EXTRA_ARGS=--node-ip=192.168.56.10" | sudo tee /etc/default/kubelet
```
**Snapshot all three VMs now: `pre-init`** — kubeadm experiments are wonderfully repeatable from here.

### 6.3 Initialize the control plane (k8s-cp only)

```bash
sudo kubeadm init \
  --apiserver-advertise-address=192.168.56.10 \
  --pod-network-cidr=10.244.0.0/16 \
  --node-name=k8s-cp
```
Success output ends with a **`kubeadm join ...` command — copy it somewhere safe.** Then:
```bash
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
kubectl get nodes    # k8s-cp: NotReady — expected, no CNI yet
```

### 6.4 Install Calico (CNI + NetworkPolicy engine)

We chose 10.244.0.0/16 for pods so it can't collide with the 192.168.56.0/24 host-only network — tell Calico the same:
```bash
curl -LO https://raw.githubusercontent.com/projectcalico/calico/v3.29.1/manifests/calico.yaml
# (check projectcalico releases for the current version tag)

# Uncomment & set the pool to match --pod-network-cidr:
sed -i 's|# - name: CALICO_IPV4POOL_CIDR|- name: CALICO_IPV4POOL_CIDR|' calico.yaml
sed -i 's|#   value: "192.168.0.0/16"|  value: "10.244.0.0/16"|' calico.yaml

kubectl apply -f calico.yaml
watch kubectl get pods -n kube-system    # wait for calico-node & coredns Running
kubectl get nodes                        # k8s-cp is now Ready
```

### 6.5 Join the workers (k8s-w1 and k8s-w2)

Paste the saved join command with sudo on each worker:
```bash
sudo kubeadm join 192.168.56.10:6443 --token <token> \
  --discovery-token-ca-cert-hash sha256:<hash>
```
Lost it? Regenerate on cp: `kubeadm token create --print-join-command`.
```bash
# on k8s-cp:
kubectl get nodes -o wide
# NAME     STATUS  ROLES          ...  INTERNAL-IP
# k8s-cp   Ready   control-plane       192.168.56.10   ← INTERNAL-IPs must be 56.x, not 10.0.2.15
# k8s-w1   Ready   <none>              192.168.56.11
# k8s-w2   Ready   <none>              192.168.56.12
```

### 6.6 Drive the cluster from your HOST machine

The API server (192.168.56.10:6443) is directly reachable — install kubectl on the host, then:
```bash
# on the host
scp k8s@192.168.56.10:~/.kube/config ~/.kube/config-vbox
export KUBECONFIG=~/.kube/config-vbox
kubectl get nodes         # you now administer a multi-node cluster from your own terminal
```

### 6.7 Smoke test — the payoff over minikube

```bash
kubectl create deployment web --image=nginx --replicas=4
kubectl expose deployment web --type=NodePort --port=80
kubectl get svc web                        # e.g. port 31234
kubectl get pods -o wide                   # pods spread across w1 AND w2
```
**Host browser → `http://192.168.56.11:31234`** ✅ — NodePorts now work directly; no port-forward pattern needed.

### 6.8 MetalLB — real LoadBalancer IPs on your host-only network

```bash
kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.14.9/config/manifests/metallb-native.yaml
kubectl wait -n metallb-system --for=condition=ready pod --all --timeout=120s

cat <<'EOF' | kubectl apply -f -
apiVersion: metallb.io/v1beta1
kind: IPAddressPool
metadata: {name: lab-pool, namespace: metallb-system}
spec:
  addresses: [192.168.56.200-192.168.56.240]
---
apiVersion: metallb.io/v1beta1
kind: L2Advertisement
metadata: {name: lab-l2, namespace: metallb-system}
spec:
  ipAddressPools: [lab-pool]
EOF

kubectl patch svc web -p '{"spec":{"type":"LoadBalancer"}}'
kubectl get svc web    # EXTERNAL-IP: 192.168.56.200 → open it in the host browser
```
**Snapshot all three: `cluster-ready`.** This is your restore point for the entire intermediate track.

---

## 7. Intermediate Track (Weeks 7–14) — Detailed Labs
> All labs target the **3-VM kubeadm cluster**, driven from your host terminal (§6.6). Snapshot all three VMs before each module.

### Module I1 — Helm: Packaging Applications
**Topics:** charts, releases, repos · values & overrides · Go templating & `_helpers.tpl` · `helm upgrade --install`, rollback · Helm vs Kustomize.

**Lab I1 — step by step:**
```bash
# 1. Consume a public chart and read what it renders
helm repo add bitnami https://charts.bitnami.com/bitnami && helm repo update
helm install demo bitnami/nginx --set service.type=LoadBalancer
kubectl get svc demo-nginx        # MetalLB hands it an IP → check in host browser
helm get manifest demo | less     # every resource Helm created
helm uninstall demo

# 2. Author your own chart for Capstone 1's app
helm create voting && cd voting
```
3. Restructure: one template per component (`vote-deploy.yaml`, `vote-svc.yaml`, `redis.yaml`, `db-statefulset.yaml` later, …), driven by values:
```yaml
# values.yaml (excerpt)
vote:
  image: {repository: dockersamples/examplevotingapp_vote, tag: latest}
  replicas: 2
  service: {type: LoadBalancer}
```
```yaml
# templates/vote-deploy.yaml (excerpt showing the pattern)
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "voting.fullname" . }}-vote
spec:
  replicas: {{ .Values.vote.replicas }}
  ...
        image: "{{ .Values.vote.image.repository }}:{{ .Values.vote.image.tag }}"
```
4. Per-environment values files (`values-dev.yaml` with 1 replica, `values-staging.yaml` with 2) and validate:
```bash
helm lint .
helm template voting . -f values-dev.yaml | kubectl apply --dry-run=server -f -
helm install voting . -n vote --create-namespace -f values-dev.yaml
helm upgrade voting . -n vote -f values-staging.yaml
helm rollback voting 1 -n vote
```
**✅ Done when:** one `helm install` stands up the entire capstone; `helm rollback` demonstrably reverts a bad values change.

---

### Module I2 — Ingress & TLS
**Topics:** controller vs resource · host/path routing · TLS termination · cert-manager issuers · hosts-file DNS for labs.

**Lab I2 — step by step:**
```bash
# 1. NGINX Ingress via Helm, exposed through MetalLB
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx && helm repo update
helm install ingress ingress-nginx/ingress-nginx -n ingress-nginx --create-namespace
kubectl get svc -n ingress-nginx ingress-ingress-nginx-controller   # e.g. EXTERNAL-IP 192.168.56.201
```
2. **Fake DNS on your HOST machine** — add to `/etc/hosts` (Windows: `C:\Windows\System32\drivers\etc\hosts`, edit as Administrator):
```
192.168.56.201  vote.local results.local
```
3. Route both frontends through one entry point:
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: voting
  namespace: vote
spec:
  ingressClassName: nginx
  rules:
  - host: vote.local
    http: {paths: [{path: /, pathType: Prefix, backend: {service: {name: voting-vote, port: {number: 80}}}}]}
  - host: results.local
    http: {paths: [{path: /, pathType: Prefix, backend: {service: {name: voting-result, port: {number: 80}}}}]}
```
Host browser → `http://vote.local` and `http://results.local` ✅ (switch the app Services back to ClusterIP — the Ingress is now the only front door).
4. **TLS with cert-manager:**
```bash
helm repo add jetstack https://charts.jetstack.io && helm repo update
helm install cert-manager jetstack/cert-manager -n cert-manager --create-namespace --set crds.enabled=true

cat <<'EOF' | kubectl apply -f -
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata: {name: selfsigned}
spec: {selfSigned: {}}
EOF
```
Add to the Ingress:
```yaml
metadata:
  annotations:
    cert-manager.io/cluster-issuer: selfsigned
spec:
  tls:
  - hosts: [vote.local, results.local]
    secretName: voting-tls
```
`https://vote.local` now serves TLS (accept the self-signed warning — inspect the certificate in the browser and find cert-manager as the issuer).
**✅ Done when:** both hosts serve HTTPS through one LoadBalancer IP and you can explain the request path: hosts-file → MetalLB IP → ingress controller pod → Service → app pod.

---

### Module I3 — Reliability & Autoscaling
**Topics:** metrics-server · HPA · QoS classes · PodDisruptionBudgets · maxSurge/maxUnavailable · draining nodes.

**Lab I3 — step by step:**
```bash
# 1. metrics-server (kubeadm needs the insecure-TLS flag for kubelet certs)
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
kubectl patch deployment metrics-server -n kube-system --type=json \
  -p='[{"op":"add","path":"/spec/template/spec/containers/0/args/-","value":"--kubelet-insecure-tls"}]'
kubectl top nodes    # numbers = working

# 2. HPA on the vote frontend (requests must be set — Helm values from I1)
kubectl autoscale deployment voting-vote -n vote --cpu-percent=50 --min=2 --max=8

# 3. Load test from INSIDE the cluster (or install k6 on the cp VM)
kubectl run load --rm -it --image=busybox:1.36 -- \
  /bin/sh -c "while true; do wget -qO- http://voting-vote.vote/ >/dev/null; done"
watch kubectl get hpa,pods -n vote    # replicas climb toward 8; Ctrl-C the load; watch scale-down (~5 min)

# 4. Disruption tolerance: PDB + a real node drain
cat <<'EOF' | kubectl apply -f -
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata: {name: vote-pdb, namespace: vote}
spec:
  minAvailable: 1
  selector: {matchLabels: {app.kubernetes.io/name: voting, component: vote}}
EOF
kubectl drain k8s-w1 --ignore-daemonsets --delete-emptydir-data
kubectl get pods -n vote -o wide      # everything rescheduled onto w2; app stayed up (test vote.local!)
kubectl uncordon k8s-w1
```
**✅ Done when:** you demo scale-out under load and a zero-downtime drain — this exact demo is Capstone 2 requirement #5.

---

### Module I4 — Stateful Workloads
**Topics:** StatefulSets vs Deployments · stable identity & ordered rollout · headless Services · volumeClaimTemplates · Jobs & CronJobs · local-path storage on kubeadm.

**Lab I4 — step by step:**
```bash
# 0. kubeadm has no default StorageClass (minikube spoiled you) — install one:
kubectl apply -f https://raw.githubusercontent.com/rancher/local-path-provisioner/master/deploy/local-path-storage.yaml
kubectl patch storageclass local-path -p '{"metadata":{"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'
```
1. **Postgres as a StatefulSet** with a headless Service:
```yaml
apiVersion: v1
kind: Service
metadata: {name: db, namespace: vote}
spec:
  clusterIP: None            # headless → per-pod DNS
  selector: {app: db}
  ports: [{port: 5432}]
---
apiVersion: apps/v1
kind: StatefulSet
metadata: {name: db, namespace: vote}
spec:
  serviceName: db
  replicas: 1
  selector: {matchLabels: {app: db}}
  template:
    metadata: {labels: {app: db}}
    spec:
      containers:
      - name: postgres
        image: postgres:16
        env:
        - name: POSTGRES_PASSWORD
          valueFrom: {secretKeyRef: {name: db-secret, key: DB_PASSWORD}}
        - {name: PGDATA, value: /var/lib/postgresql/data/pgdata}
        volumeMounts: [{name: data, mountPath: /var/lib/postgresql/data}]
  volumeClaimTemplates:
  - metadata: {name: data}
    spec:
      accessModes: [ReadWriteOnce]
      resources: {requests: {storage: 2Gi}}
```
2. **See what "stateful" buys you** — scale a 3-replica Redis StatefulSet and inspect:
```bash
kubectl get pods -n vote -l app=redis     # redis-0, redis-1, redis-2 — names, not hashes
kubectl get pvc -n vote                   # data-redis-0/1/2 — one claim per pod
kubectl run tmp -n vote --rm -it --image=busybox:1.36 -- nslookup redis-0.redis.vote.svc.cluster.local
kubectl delete pod redis-1 -n vote        # returns as redis-1 with the SAME PVC
```
3. **CronJob backup:**
```yaml
apiVersion: batch/v1
kind: CronJob
metadata: {name: pg-backup, namespace: vote}
spec:
  schedule: "0 2 * * *"
  jobTemplate:
    spec:
      template:
        spec:
          restartPolicy: OnFailure
          containers:
          - name: dump
            image: postgres:16
            command: ["/bin/sh","-c","pg_dump -h db -U postgres postgres > /backup/db-$(date +%F).sql"]
            env:
            - name: PGPASSWORD
              valueFrom: {secretKeyRef: {name: db-secret, key: DB_PASSWORD}}
            volumeMounts: [{name: backup, mountPath: /backup}]
          volumes:
          - name: backup
            persistentVolumeClaim: {claimName: backup-pvc}
```
Trigger it now instead of waiting for 02:00: `kubectl create job -n vote --from=cronjob/pg-backup backup-now`, then exec in and `ls /backup`.
**✅ Done when:** you can demo per-pod identity + storage surviving deletion, and articulate when you'd still run the database *outside* the cluster.

---

### Module I5 — Security Fundamentals
**Topics:** authn vs authz · RBAC · ServiceAccounts · SecurityContext · NetworkPolicies (Calico enforces them — this is why we chose it) · Trivy scanning.

**Lab I5 — step by step:**
1. **Least-privilege human access:**
```bash
kubectl create serviceaccount developer -n vote
cat <<'EOF' | kubectl apply -f -
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata: {name: pod-reader, namespace: vote}
rules:
- apiGroups: [""]
  resources: [pods, pods/log]
  verbs: [get, list, watch]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata: {name: developer-pod-reader, namespace: vote}
subjects: [{kind: ServiceAccount, name: developer, namespace: vote}]
roleRef: {kind: Role, name: pod-reader, apiGroup: rbac.authorization.k8s.io}
EOF
# Prove the boundary:
kubectl auth can-i list pods  -n vote --as=system:serviceaccount:vote:developer     # yes
kubectl auth can-i delete pods -n vote --as=system:serviceaccount:vote:developer    # no
kubectl auth can-i list pods  -n default --as=system:serviceaccount:vote:developer  # no
```
2. **Default-deny the namespace, then allow only legitimate flows:**
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata: {name: default-deny, namespace: vote}
spec:
  podSelector: {}
  policyTypes: [Ingress, Egress]
```
The app breaks (good!). Now write allow policies: vote→redis:6379, worker→redis & →db:5432, result→db:5432, ingress-controller→frontends, all→kube-dns:53 (UDP). Verify the wall:
```bash
kubectl run attacker -n vote --rm -it --image=busybox:1.36 -- nc -zv -w 2 db 5432   # must TIME OUT
```
3. **Scan & harden:**
```bash
# Trivy on the cp VM:
sudo apt-get install -y wget apt-transport-https gnupg
wget -qO- https://aquasecurity.github.io/trivy-repo/deb/public.key | sudo gpg --dearmor -o /usr/share/keyrings/trivy.gpg
echo "deb [signed-by=/usr/share/keyrings/trivy.gpg] https://aquasecurity.github.io/trivy-repo/deb generic main" | sudo tee /etc/apt/sources.list.d/trivy.list
sudo apt-get update && sudo apt-get install -y trivy
trivy image dockersamples/examplevotingapp_vote   # triage HIGH/CRITICAL findings
```
Harden one Deployment and keep it working:
```yaml
      securityContext:
        runAsNonRoot: true
        runAsUser: 10001
      containers:
      - ...
        securityContext:
          allowPrivilegeEscalation: false
          readOnlyRootFilesystem: true
          capabilities: {drop: [ALL]}
```
**✅ Done when:** the attacker pod can't reach the DB, `auth can-i` proves RBAC boundaries, and you fixed ≥1 HIGH finding.

---

### Module I6 — Observability: Metrics, Dashboards & Logs
**Topics:** Prometheus scraping & PromQL basics · Grafana · kube-state-metrics · Loki · Alertmanager · golden signals.

**Lab I6 — step by step:**
```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add grafana https://grafana.github.io/helm-charts && helm repo update

helm install monitoring prometheus-community/kube-prometheus-stack -n monitoring --create-namespace \
  --set grafana.service.type=LoadBalancer
helm install loki grafana/loki-stack -n monitoring --set promtail.enabled=true

kubectl get svc -n monitoring monitoring-grafana     # e.g. 192.168.56.202
kubectl get secret -n monitoring monitoring-grafana -o jsonpath='{.data.admin-password}' | base64 -d; echo
```
Host browser → `http://192.168.56.202` (user `admin`, password above).
1. Import community dashboard **ID 15757** (Kubernetes views) → watch your own cluster.
2. Add Loki as a data source (`http://loki:3100`) → Explore → `{namespace="vote"}` → your app's logs.
3. **Build your own capstone dashboard:** pod restarts (`kube_pod_container_status_restarts_total`), CPU vs limits, replica counts, HPA activity. Re-run the I3 load test and watch it live.
4. **One real alert:**
```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: vote-alerts
  namespace: monitoring
  labels: {release: monitoring}
spec:
  groups:
  - name: vote
    rules:
    - alert: PodRestartLoop
      expr: increase(kube_pod_container_status_restarts_total{namespace="vote"}[10m]) > 3
      for: 1m
      labels: {severity: warning}
      annotations: {summary: "Pod {{ $labels.pod }} is restart-looping"}
```
Fire it on purpose (deploy a crashing image) → watch it turn red in Prometheus → Alerts.
**✅ Done when:** given a deliberately broken deploy, you find the cause using ONLY Grafana + Loki — kubectl stays closed.

---

### Module I7 — CI/CD & GitOps
**Topics:** build→test→scan→push pipelines · GitHub Actions · image tagging · GitOps principles · Argo CD sync, drift & rollback.

**Lab I7 — step by step:**
1. **App repo** (fork/copy the voting app or your own service) — `.github/workflows/ci.yaml`:
```yaml
name: build-scan-push
on: {push: {branches: [main]}}
permissions: {contents: read, packages: write}
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - uses: docker/login-action@v3
      with: {registry: ghcr.io, username: ${{ github.actor }}, password: ${{ secrets.GITHUB_TOKEN }}}
    - uses: docker/build-push-action@v6
      with:
        context: ./vote
        push: true
        tags: ghcr.io/${{ github.repository }}/... (9 KB left)