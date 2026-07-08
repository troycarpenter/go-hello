# go-hello

A simple Go application used to validate an end-to-end CI/CD deployment workflow.

The application is intentionally minimal. The purpose of this project is to test the infrastructure and automation required to move an application from source code into a running Kubernetes environment.

## Pipeline Flow

Source Change
|
v
Jenkins Pipeline
|
v
Go Build
|
v
Docker Image Build
|
v
Harbor Container Registry
|
v
Kubernetes Deployment

## Components

- Go application
- Jenkins pipeline automation
- Docker container image
- Harbor image registry
- Kubernetes deployment manifests
- k3s cluster deployment target

## Purpose

This project was created to validate:

- CI/CD automation workflows
- Container image lifecycle management
- Kubernetes deployment automation
- Integration between build systems and infrastructure platforms

The application is intentionally simple because the focus is the deployment process rather than application functionality.
