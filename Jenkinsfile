pipeline {
  agent none

  environment {
    REGISTRY = "harbor.carpenter.cx"
    IMAGE = "library/go-hello"
    TAG = "${BUILD_NUMBER}"
  }

  stages {

    stage('Build and Deploy') {

      agent {
        kubernetes {
          label "go-kaniko-${BUILD_NUMBER}"
          defaultContainer 'golang'

          yaml """
apiVersion: v1
kind: Pod
spec:
  containers:

  - name: golang
    image: golang:1.25
    command: ["sleep"]
    args: ["infinity"]

  - name: kaniko
    image: gcr.io/kaniko-project/executor:v1.23.2
    command: ["sleep"]
    args: ["infinity"]
    volumeMounts:
      - name: docker-config
        mountPath: /kaniko/.docker

  - name: kubectl
    image: bitnami/kubectl:latest
    command: ["sleep"]
    args: ["infinity"]

  volumes:
    - name: docker-config
      secret:
        secretName: harbor-regcred

  restartPolicy: Never
"""
        }
      }

      stages {

        stage('Checkout') {
          steps {
            checkout scm
          }
        }

        stage('Build') {
          steps {
            container('golang') {
              sh 'go version'
            }
          }
        }

        stage('Build Image') {
          steps {
            container('kaniko') {
              sh """
              /kaniko/executor \
                --context=${WORKSPACE} \
                --dockerfile=${WORKSPACE}/Dockerfile \
                --destination=${REGISTRY}/${IMAGE}:${TAG} \
                --cleanup
              """
            }
          }
        }

        stage('Deploy') {
          steps {
            container('kubectl') {
              sh """
              kubectl apply -f k8s/deployment.yaml
              kubectl set image deployment/go-hello \
                go-hello=${REGISTRY}/${IMAGE}:${TAG}
              """
            }
          }
        }
      }
    }
  }
}
