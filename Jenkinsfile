node { 
    stage('Dependency') {
	sh "make restore"
    }
    stage('CodeGen') {
        sh "make codegen"
    }
    stage('Build') {
        sh "make all"
    }
    stage('Lint') {
        sh "make lint"
    }
}
