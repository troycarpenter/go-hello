pipeline {
  agent none

  environment {
    REGISTRY = "harbor.carpenter.cx"
    IMAGE    = "library/go-hello"
    TAG      = "${BUILD_NUMBER}"
    FULL_IMAGE = "${REGISTRY}/${IMAGE}:${BUILD_NUMBER}"
  }

  stages {

    stage('Checkout') {
      agent any
      steps {
        checkout scm
      }
    }

    stage('Build Image (Kaniko Job)') {
      agent {
        kubernetes {
          label "kaniko-${BUILD_NUMBER}"
          defaultContainer 'jnlp'

          yaml """
apiVersion: v1
kind: Pod
spec:
  restartPolicy: Never

  containers:
  - name: kubectl
    image: bitnami/kubectl:latest
    command: ['cat']
    tty: true
"""
        }
      }

      steps {
        container('kubectl') {
          sh """
          sed \
            -e 's|__IMAGE__|${FULL_IMAGE}|g' \
            -e 's|__BUILD_NUMBER__|${BUILD_NUMBER}|g' \
            k8s/kaniko-job.yaml | kubectl apply -f -

          kubectl wait --for=condition=complete job/go-hello-build --timeout=600s
          kubectl logs job/go-hello-build
          """
        }
      }
    }

    stage('Deploy') {
      agent {
        kubernetes {
          label "deploy-${BUILD_NUMBER}"
          defaultContainer 'kubectl'

          yaml """
apiVersion: v1
kind: Pod
spec:
  restartPolicy: Never

  containers:
  - name: kubectl
    image: bitnami/kubectl:1.30
    command: ['cat']
    tty: true
"""
        }
      }

      steps {
        container('kubectl') {
          sh """
          kubectl apply -f k8s/deployment.yaml

          kubectl set image deployment/go-hello \
            go-hello=${FULL_IMAGE} \
            -n default

          kubectl rollout status deployment/go-hello -n default
          """
        }
      }
    }
  }
}
