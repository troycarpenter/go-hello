pipeline {
    agent any

    environment {
        // ── Update these for your environment ──────────────────────────────
        HARBOR_REGISTRY  = "harbor.carpenter.cx"          // Harbor hostname/IP
        HARBOR_PROJECT   = "library"                       // Harbor project name
        APP_NAME         = "go-hello"                       // Image name
        DEPLOY_NAMESPACE = "default"                         // k8s namespace
        DEPLOY_NAME      = "go-hello"                       // k8s Deployment name
        // ───────────────────────────────────────────────────────────────────

        IMAGE_TAG        = "${HARBOR_REGISTRY}/${HARBOR_PROJECT}/${APP_NAME}"
        GIT_COMMIT_SHORT = "${GIT_COMMIT[0..7]}"
    }

    options {
        buildDiscarder(logRotator(numToKeepStr: '10'))
        timeout(time: 30, unit: 'MINUTES')
        disableConcurrentBuilds()
    }

    triggers {
        // Poll GitHub every minute as a fallback; webhook is preferred (see Step 3)
        pollSCM('* * * * *')
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
            agent {
                docker {
                    image 'golang:1.22-alpine'
                    args  '-v /tmp/go-cache:/go/pkg/mod'   // cache Go modules
                    reuseNode true
                }
            }
            steps {
                sh '''
                    go env -w GOFLAGS=-mod=mod
                    go vet ./...
                    go test -v -race -coverprofile=coverage.out ./...
                '''
            }
            post {
                always {
                    // Publish test results if you use gotestsum/junit output
                    // junit 'reports/*.xml'
                    archiveArtifacts artifacts: 'coverage.out', allowEmptyArchive: true
                }
            }
        }

        stage('Build Docker Image') {
            steps {
                script {
                    docker.build(
                        "${IMAGE_TAG}:${GIT_COMMIT_SHORT}",
                        "--label git-commit=${GIT_COMMIT_SHORT} " +
                        "--label build-date=${new Date().format('yyyy-MM-dd')} " +
                        "--label branch=${GIT_BRANCH_CLEAN} " +
                        "."
                    )
                }
            }
        }

        stage('Push to Harbor') {
            steps {
                withCredentials([
                    usernamePassword(
                        credentialsId: 'harbor-credentials',
                        usernameVariable: 'HARBOR_USER',
                        passwordVariable: 'HARBOR_PASS'
                    )
                ]) {
                    script {
                        sh "echo '${HARBOR_PASS}' | docker login ${HARBOR_REGISTRY} -u '${HARBOR_USER}' --password-stdin"

                        // Tag and push with commit SHA
                        sh "docker push ${IMAGE_TAG}:${GIT_COMMIT_SHORT}"

                        // Also tag as 'latest' on main/master branch
                        if (env.GIT_BRANCH_CLEAN in ['main', 'master', 'origin-main', 'origin-master']) {
                            sh "docker tag  ${IMAGE_TAG}:${GIT_COMMIT_SHORT} ${IMAGE_TAG}:latest"
                            sh "docker push ${IMAGE_TAG}:latest"
                        }
                    }
                }
            }
            post {
                always {
                    sh "docker logout ${HARBOR_REGISTRY} || true"
                }
            }
        }

        stage('Deploy to k3s') {
            // Remove this stage if you prefer GitOps (ArgoCD/Flux) for deployment
            when {
                anyOf {
                    branch 'main'
                    branch 'master'
                }
            }
            steps {
                withCredentials([
                    file(credentialsId: 'kubeconfig', variable: 'KUBECONFIG')
                ]) {
                    sh """
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

    post {
        success {
            echo "✅ Pipeline succeeded. Image: ${IMAGE_TAG}:${GIT_COMMIT_SHORT}"
        }
        failure {
            echo "❌ Pipeline failed. Check logs above."
            // Add email/Slack notification here if desired
        }
        cleanup {
            // Remove local images to keep the Jenkins node disk clean
            sh "docker rmi ${IMAGE_TAG}:${GIT_COMMIT_SHORT} || true"
            sh "docker rmi ${IMAGE_TAG}:latest || true"
            cleanWs()
        }
    }
}
