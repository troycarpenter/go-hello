pipeline {
  agent any

  environment {
    REGISTRY = "harbor.carpenter.cx"
    IMAGE    = "library/go-hello"
    TAG      = "${BUILD_NUMBER}"
    FULL_IMAGE = "${REGISTRY}/${IMAGE}:${TAG}"
  }

  stages {

    stage('Checkout') {
      steps {
        checkout scm
      }
    }

    stage('Build Image (Kaniko Job)') {
      steps {
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

    stage('Deploy') {
      steps {
        sh """
        kubectl apply -f k8s/deployment.yaml

        kubectl set image deployment/go-hello \
          go-hello=${FULL_IMAGE} \
          -n default
        """
      }
    }
  }
}
