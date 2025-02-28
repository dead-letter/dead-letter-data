# Microservice Template
A template with an automatic workflow for building & publishing Docker Images

## Usage
1. Configure the personal access token under the CR_PAT variable in the GitHub Actions Secrets
2. On push, the workflow will run and 
    - Build the Docker image
    - Publish the Docker image under USERNAME/REPO_NAME:COMMIT_SHA to the GHCR (not Docker hub)
