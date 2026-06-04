pipeline {
  agent none

  environment {
    REGISTRY = "harbor.carpenter.cx"
    IMAGE = "library/go-hello"
    TAG = "${BUILD_NUMBER}"
  }

  stages {

    stage('Checkout') {
      agent any
      steps {
        checkout scm
      }
    }

    stage('Build Image (Kaniko)') {
      agent {
        kubernetes {
          label "kaniko-${BUILD_NUMBER}"
          defaultContainer 'kaniko'

          yaml """
apiVersion: v1
kind: Pod
spec:
  restartPolicy: Never
  containers:

  - name: kaniko
    image: gcr.io/kaniko-project/executor:v1.23.2
    command:
    - /kaniko/executor
    args:
    - --context=/workspace
    - --dockerfile=/workspace/Dockerfile
    - --destination=harbor.carpenter.cx/library/go-hello:${BUILD_NUMBER}
    - --cleanup
    volumeMounts:
    - name: docker-config
      mountPath: /kaniko/.docker

  volumes:
  - name: docker-config
    secret:
      secretName: harbor-regcred
"""
        }
      }

      steps {
        container('kaniko') {
          sh """
            /kaniko/executor \
              --context=$WORKSPACE \
              --dockerfile=$WORKSPACE/Dockerfile \
              --destination=${REGISTRY}/${IMAGE}:${TAG} \
              --cleanup
          """
        }
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
    command: ["sh", "-c", "cat"]
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
  }
}
