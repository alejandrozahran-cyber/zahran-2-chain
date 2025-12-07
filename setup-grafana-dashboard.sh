#!/bin/bash

cd ~/zahran-2-chain

echo "Setting up Grafana dashboard..."

# Create directories
mkdir -p docker/grafana/provisioning/dashboards
mkdir -p docker/grafana/provisioning/datasources

# Create provisioning config
cat > docker/grafana/provisioning/dashboards/default. yaml << 'YAML_END'
apiVersion: 1
providers:
  - name: 'NUSA Chain'
    orgId: 1
    folder: ''
    type: file
    disableDeletion: false
    updateIntervalSeconds: 10
    allowUiUpdates: true
    options:
      path: /etc/grafana/provisioning/dashboards
      foldersFromFilesStructure: true
YAML_END

# Create dashboard JSON (simplified)
cat > docker/grafana/provisioning/dashboards/nusa-dashboard.json << 'JSON_END'
{
  "title": "NUSA Chain Metrics",
  "panels": [],
  "schemaVersion": 16,
  "version": 0,
  "uid": "nusa-main"
}
JSON_END

echo "✅ Files created!"
echo "Now wait 10 seconds and refresh Grafana..."

