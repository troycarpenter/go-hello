pipeline {
    agent {
        kubernetes {
            yaml """
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: golang
    image: golang:1.22-alpine
    command: ['sleep', '99d']
    resources:
      requests:
        memory: "512Mi"
        cpu: "500m"
  - name: docker
    image: docker:24-dind
    securityContext:
      privileged: true
    env:
    - name: DOCKER_TLS_CERTDIR
      value: ""
    resources:
      requests:
        memory: "512Mi"
        cpu: "500m"
    volumeMounts:
    - name: docker-storage
      mountPath: /var/lib/docker
  volumes:
  - name: docker-storage
    emptyDir: {}
"""
        }
    }
 
    environment {
        HARBOR_REGISTRY  = "harbor.carpenter.cx"
        HARBOR_PROJECT   = "library"
        APP_NAME         = "go-hello"
        DEPLOY_NAMESPACE = "default"
        DEPLOY_NAME      = "go-hello"
 
        IMAGE_TAG        = "${HARBOR_REGISTRY}/${HARBOR_PROJECT}/${APP_NAME}"
    }
 
    options {
        buildDiscarder(logRotator(numToKeepStr: '10'))
        timeout(time: 30, unit: 'MINUTES')
        disableConcurrentBuilds()
    }
 
    stages {
 
        stage('Checkout') {
            steps {
                checkout scm
                script {
                    env.GIT_COMMIT_SHORT = sh(
                        script: 'git rev-parse --short HEAD',
                        returnStdout: true
                    ).trim()
                    env.GIT_BRANCH_CLEAN = env.GIT_BRANCH?.replaceAll('/', '-') ?: 'unknown'
                    echo "Building commit: ${env.GIT_COMMIT_SHORT} on branch: ${env.GIT_BRANCH_CLEAN}"
                }
            }
        }
 
        stage('Test') {
            steps {
                container('golang') {
                    sh '''
                        go env -w GOFLAGS=-mod=mod
                        go vet ./...
                        go test -v -coverprofile=coverage.out ./...
                    '''
                }
            }
            post {
                always {
                    archiveArtifacts artifacts: 'coverage.out', allowEmptyArchive: true
                }
            }
        }
 
        stage('Build Docker Image') {
            steps {
                container('docker') {
                    sh "docker build -t ${IMAGE_TAG}:${GIT_COMMIT_SHORT} ."
                }
            }
        }
 
        stage('Push to Harbor') {
            steps {
                container('docker') {
                    withCredentials([
                        usernamePassword(
                            credentialsId: 'harbor-credentials',
                            usernameVariable: 'HARBOR_USER',
                            passwordVariable: 'HARBOR_PASS'
                        )
                    ]) {
                        sh "echo '${HARBOR_PASS}' | docker login ${HARBOR_REGISTRY} -u '${HARBOR_USER}' --password-stdin"
                        sh "docker push ${IMAGE_TAG}:${GIT_COMMIT_SHORT}"
 
                        script {
                            if (env.GIT_BRANCH_CLEAN in ['main', 'master', 'origin-main', 'origin-master']) {
                                sh "docker tag ${IMAGE_TAG}:${GIT_COMMIT_SHORT} ${IMAGE_TAG}:latest"
                                sh "docker push ${IMAGE_TAG}:latest"
                            }
                        }
                    }
                }
            }
            post {
                always {
                    container('docker') {
                        sh "docker logout ${HARBOR_REGISTRY} || true"
                    }
                }
            }
        }
 
        stage('Deploy to k3s') {
            when {
                expression {
                     env.GIT_BRANCH_CLEAN in ['main', 'master', 'origin-main', 'origin-master'] 
                }
            }
            steps {
                withCredentials([
                    file(credentialsId: 'kubeconfig', variable: 'KUBECONFIG')
                ]) {
                    container('golang') {
                        sh """
                            apk add --no-cache curl
                            curl -LO https://dl.k8s.io/release/v1.29.0/bin/linux/amd64/kubectl
                            chmod +x kubectl && mv kubectl /usr/local/bin/
 
                            kubectl --kubeconfig=\$KUBECONFIG \
                                set image deployment/${DEPLOY_NAME} \
                                ${APP_NAME}=${IMAGE_TAG}:${GIT_COMMIT_SHORT} \
                                -n ${DEPLOY_NAMESPACE}
 
                            kubectl --kubeconfig=\$KUBECONFIG \
                                rollout status deployment/${DEPLOY_NAME} \
                                -n ${DEPLOY_NAMESPACE} \
                                --timeout=120s
                        """
                    }
                }
            }
        }
    }
 
    post {
        success {
            echo "✅ Pipeline succeeded. Image: ${IMAGE_TAG}:${GIT_COMMIT_SHORT}"
        }
        failure {
            echo "❌ Pipeline failed. Check logs above."
        }
        cleanup {
            deleteDir()
        }
    }
}
