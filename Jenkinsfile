#!/usr/bin/env groovy

properties([
        disableConcurrentBuilds(),
])

node {
    def goHome = "${tool 'Go 1.13'}"
    def go = "${goHome}/bin/go"
    def dockerImageName = "kong-federated-ingress"
    def version = env.BUILD_ID
    def toolsNexus = "toolsnexus.marchex.com"

    timestamps {
        stage('Checkout') {
            checkout scm
            sh 'git branch --set-upstream-to=origin/master master'
        }
        stage('Test') {
//             sh "${go} test ./..."
        }
        stage('Compile') {
            sh "mkdir bin"
            sh "${go} build -o bin/ ./cmd/k8s-kong-federated-ingress/"
        }
        stage('Build and Publish Docker Image') {
            def image
            docker.withRegistry("https://${toolsNexus}:5000", 'oce-build-automation') {
                image = docker.build("${dockerImageName}:${version}")
            }
            docker.withRegistry("https://${toolsNexus}:5001", 'oce-build-automation') {
                image.push()
                image.push("latest")
            }
            sshagent(credentials: ['oce-build-automation']) {
                sh "git tag -a 'docker-image-${version}' -m '[jenkins] Tagging release version ${version}'"
                sh "git push --tags"
            }
            currentBuild.description = "Docker Image: ${dockerImageName}:${version}"
        }
    }
}
