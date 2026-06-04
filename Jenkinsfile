pipeline {
  agent any

  environment {
    REGISTRY = "harbor.carpenter.cx"
    IMAGE = "library/go-hello"
    TAG = "${BUILD_NUMBER}"
  }

  stages {

    stage('Checkout') {
      steps {
        checkout scm
      }
    }

    stage('Build Image (K8s Job)') {
      steps {
        sh """
        envsubst < k8s/kaniko-job.yaml | kubectl apply -f -
        kubectl wait --for=condition=complete job/go-hello-build --timeout=300s
        kubectl logs job/go-hello-build
        """
      }
    }

stage('Deploy') {
  agent {
    kubernetes {
      label "kubectl-${BUILD_NUMBER}"
      defaultContainer 'kubectl'

      yaml """
apiVersion: v1
kind: Pod
spec:
  restartPolicy: Never
  containers:
  - name: kubectl
    image: bitnami/kubectl:1.30
    command: ["cat"]
    tty: true
"""
    }
  }

  steps {
    container('kubectl') {
      sh """
      kubectl apply -f k8s/deployment.yaml

      kubectl set image deployment/go-hello \
        go-hello=${REGISTRY}/${IMAGE}:${TAG} \
        -n default
      """
    }
  }
}
