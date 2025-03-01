# Microservice Template
A template with an automatic workflow for building & publishing Docker Images

## Usage
1. Configure the personal access token under the CR_PAT variable in the GitHub Actions Secrets
2. On push, the workflow will run and 
    - Build the Docker image
    - Publish the Docker image under USERNAME/REPO_NAME:COMMIT_SHA to the GHCR (not Docker hub)

## How to Update From Template

`git remote add template git@github.com:dead-letter/microservice-template.git`

`git fetch --all`

`git merge template/main --allow-unrelated-histories`

[Stack Overflow](https://stackoverflow.com/questions/56577184/github-pull-changes-from-a-template-repository)
