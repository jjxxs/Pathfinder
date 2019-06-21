pipeline {
    agent none
    stages {
        stage('Build') {
            agent {
                docker { image 'obraun/vss-jenkins' }
            }
            steps {
                sh 'go build main.go'
            }
        }
        stage('Test') {
            agent {
                docker { image 'obraun/vss-jenkins' }
            }
            steps {
                sh 'echo go test -v'
                sh 'echo go test -bench=.'
            }
        }
        stage('Lint') {
            agent {
                docker { image 'obraun/vss-jenkins' }
            }   
            steps {
                sh 'cd algorithm && golangci-lint run --deadline 20m --enable-all'
                sh 'cd problem && golangci-lint run --deadline 20m --enable-all'
            }
        }
        stage('Docker') {
            agent any
            steps {
                sh "docker-build-and-push -b ${BRANCH_NAME} -s webapp -f webapp.dockerfile"
                sh "docker-build-and-push -b ${BRANCH_NAME} -s solver -f solver.dockerfile"
            }
        }
    }
}
