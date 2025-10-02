#!/usr/bin/env bash
set -euo pipefail

ROOT="lostfound"

mkdir -p "${ROOT}"/services/{api-gateway/{cmd,internal,web/templates,web/static},auth/{cmd,internal},items/{cmd,internal,migrations},media/{cmd,internal,migrations},chat/{cmd,internal},match/{cmd,internal}}
mkdir -p "${ROOT}"/{libs,infra/{k8s},docs,tooling/scripts,configs,web/templates,web/static,.github/workflows}

# Top-level files
touch "${ROOT}/README.md" "${ROOT}/Makefile" "${ROOT}/docker-compose.yml"

# Per-service starter files
for svc in api-gateway auth items media chat match; do
  svcdir="${ROOT}/services/${svc}"
  touch "${svcdir}/Dockerfile" "${svcdir}/README.md"
  mkdir -p "${svcdir}/cmd"
  cat > "${svcdir}/cmd/main.go" <<'GO' 
package main

import "fmt"

func main() {
    fmt.Println("placeholder for SERVICE_NAME - replace SERVICE_NAME in the file")
}
GO
  # replace placeholder SERVICE_NAME with actual service name
  sed -i "s/SERVICE_NAME/${svc}/g" "${svcdir}/cmd/main.go"
done

# Example migrations folders and placeholder
touch "${ROOT}/services/items/migrations/README.md" "${ROOT}/services/media/migrations/README.md"

# Example config
cat > "${ROOT}/configs/example.env" <<'ENV'
# example env values
DATABASE_URL=postgres://user:pass@localhost:5432/<service_db>?sslmode=disable
MINIO_ENDPOINT=http://localhost:9000
MINIO_ACCESS_KEY=minio
MINIO_SECRET_KEY=minio123
JWT_SECRET=replace-me
ENV

echo "Project skeleton created at ./${ROOT}"
