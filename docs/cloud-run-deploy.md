# Deploying to Google Cloud Run

This app runs on Cloud Run with zero application changes — the existing Dockerfile, distroless image, and `PORT` env var are all Cloud Run needs.

## Prerequisites

- Google Cloud account with billing enabled
- `gcloud` CLI installed and authenticated
- Docker installed locally (only if building locally)

## 1. One-Time GCP Setup

```bash
export GCP_PROJECT="your-project-id"
export GCP_REGION="us-west1"

gcloud auth login
gcloud config set project $GCP_PROJECT

gcloud services enable \
  run.googleapis.com \
  artifactregistry.googleapis.com \
  cloudbuild.googleapis.com
```

## 2. Create Artifact Registry Repository

```bash
gcloud artifacts repositories create stylesheets \
  --repository-format=docker \
  --location=$GCP_REGION \
  --description="Stylesheets app container images"

gcloud auth configure-docker ${GCP_REGION}-docker.pkg.dev
```

## 3. Build the Container Image

### Option A: Remote build with Cloud Build (recommended)

No local Docker needed — Cloud Build reads the Dockerfile and builds on Google's infrastructure.

```bash
gcloud builds submit \
  --tag ${GCP_REGION}-docker.pkg.dev/${GCP_PROJECT}/stylesheets/app:latest \
  .
```

### Option B: Build locally and push

```bash
docker build -t ${GCP_REGION}-docker.pkg.dev/${GCP_PROJECT}/stylesheets/app:latest .
docker push ${GCP_REGION}-docker.pkg.dev/${GCP_PROJECT}/stylesheets/app:latest
```

## 4. Deploy to Cloud Run

```bash
gcloud run deploy stylesheets \
  --image ${GCP_REGION}-docker.pkg.dev/${GCP_PROJECT}/stylesheets/app:latest \
  --region $GCP_REGION \
  --platform managed \
  --allow-unauthenticated \
  --port 8080 \
  --memory 128Mi \
  --cpu 1 \
  --min-instances 0 \
  --max-instances 3 \
  --concurrency 80
```

| Flag | Why |
|---|---|
| `--allow-unauthenticated` | Public website, no auth needed |
| `--port 8080` | Matches Dockerfile EXPOSE and app default |
| `--memory 128Mi` | App is tiny, no need for more |
| `--min-instances 0` | Scale to zero when idle (saves cost) |
| `--max-instances 3` | Cap scaling to control costs |
| `--concurrency 80` | Go handles concurrent requests well |

The output includes your service URL: `https://stylesheets-XXXXX-uc.a.run.app`

## 5. Verify

```bash
# Check service status
gcloud run services describe stylesheets --region $GCP_REGION

# Read logs
gcloud run services logs read stylesheets --region $GCP_REGION
```

Visit the `.run.app` URL and confirm all guides load and HTMX interactions work over HTTPS.

## Continuous Deployment (Optional)

A `cloudbuild.yaml` is included in the repo root. To auto-deploy on every push to `main`:

```bash
gcloud builds triggers create github \
  --repo-name=stylesheets \
  --repo-owner=johnfarrell \
  --branch-pattern='^main$' \
  --build-config=cloudbuild.yaml
```

This builds a new image tagged with the commit SHA and deploys it to Cloud Run automatically.

## Custom Domain (Optional)

```bash
gcloud run domain-mappings create \
  --service stylesheets \
  --domain your-domain.com \
  --region $GCP_REGION
```

Add the DNS records from the output to your domain registrar. Cloud Run provides automatic TLS.

## Security Hardening (Optional)

**Dedicated service account** (instead of default compute SA):

```bash
gcloud iam service-accounts create stylesheets-sa \
  --display-name="Stylesheets Cloud Run SA"

gcloud run services update stylesheets \
  --service-account=stylesheets-sa@${GCP_PROJECT}.iam.gserviceaccount.com
```

**Restrict access** (if you add auth later):

```bash
gcloud run services update stylesheets --no-allow-unauthenticated
```

## Cost

With Cloud Run's free tier (2M requests/month, 360k vCPU-seconds), a personal style guide site will almost certainly be free. Scale-to-zero means no cost when nobody is visiting.
