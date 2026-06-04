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

          defaultContainer 'jnlp'

          yaml """
apiVersion: v1
kind: Pod
spec:
  containers:

  - name: golang
    image: golang:1.25
    command:
    - cat
    tty: true

  - name: kaniko
    image: gcr.io/kaniko-project/executor:v1.23.2
    tty: true
    volumeMounts:
    - name: docker-config
      mountPath: /kaniko/.docker
        
  - name: kubectl
    image: bitnami/kubectl:latest
    command:
    - cat
    tty: true

  restartPolicy: Never

  volumes:
  - name: docker-config
    secret:
      secretName: harbor-regcred
"""
        }
      }

      stages {

        stage('Checkout') {
          steps {
            checkout scm
          }
        }

        stage('Build with Kaniko') {
          steps {
            container('kaniko') {
              sh """
              /kaniko/executor \
              --context=${WORKSPACE} \
              --dockerfile=${WORKSPACE}/Dockerfile \
              --destination=harbor.carpenter.cx/library/go-hello:${BUILD_NUMBER} \
              --cleanup
              """
            }
          }
        }

        stage('Deploy to k3s') {
          steps {
            container('kubectl') {
              sh """
              kubectl apply -f k8s/deployment.yaml
              kubectl set image deployment/go-hello \
                go-hello=${REGISTRY}/${IMAGE}:${TAG} -n default
              """
            }
          }
        }
      }
    }
  }
}
