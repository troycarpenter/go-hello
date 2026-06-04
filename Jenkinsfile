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

    stage('Build with Kaniko') {
      steps {
        container('kaniko') {
          sh """
          /kaniko/executor \
            --context ${WORKSPACE} \
            --dockerfile Dockerfile \
            --destination ${REGISTRY}/${IMAGE}:${TAG}
          """
        }
      }
    }

    stage('Deploy to k3s') {
      steps {
        sh """
        kubectl apply -f k8s/deployment.yaml
        kubectl set image deployment/go-hello \
          go-hello=${REGISTRY}/${IMAGE}:${TAG} -n default
        """
      }
    }
  }
}
