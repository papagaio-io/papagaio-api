#!groovy

def appName = "papagaio-api"
def targetMap = [ master: 'dev', stable: 'uat', va: 'va']
def branch = env.BRANCH_NAME
def target = targetMap[branch]
def label = appName + "-${UUID.randomUUID().toString()}"

podTemplate(
  label: 'node',
  containers: [
    containerTemplate(
      image: 'docker',
      name: 'docker',
      command: 'cat',
      ttyEnabled: true)
    ],
  volumes: [
    hostPathVolume(hostPath: '/var/run/docker.sock', mountPath: '/var/run/docker.sock')])

  {
    podTemplate(
      label: 'go',
      containers: [
        containerTemplate(
          image: 'golang:1.15.8',
          name: 'go',
          command: 'cat',
          ttyEnabled: true),
        containerTemplate(
          image: 'sonarsource/sonar-scanner-cli',
          name: 'sonar-scanner',
          command: 'cat',
          ttyEnabled: true),
        containerTemplate(
          image: 'docker',
          name: 'docker',
          command: 'cat',
          ttyEnabled: true)
        ],
      volumes: [
        hostPathVolume(hostPath: '/usr/bin/kubectl', mountPath: '/usr/bin/kubectl'),
        hostPathVolume(hostPath: '/var/run/docker.sock', mountPath: '/var/run/docker.sock'),
        secretVolume(mountPath: '/etc/kubernetes', secretName: 'cluster-admin')
      ]
    )

    {
      node("go") {
        def version

        stage ('Checkout') {
          cleanWs()
          scmVars = git credentialsId: 'git', url: "git@wecode.sorint.it:opensource/papagaio-api.git", branch: "${branch}"
        }


        stage ('Version') {
          if (branch == 'master') {
              version = 'latest'
          } else {
              version = sh (script: 'egrep -o "[0-9]+\\.[0-9]+\\.[0-9]+" main.go', returnStdout: true).trim()
           }
        }

        container('go') {
          withEnv(["GOROOT=/usr/local/go", "GOPATH=${PWD}:${PWD}/src/wecode.sorint.it/opensource/papagaio-api"]) {
            env.PATH="${GOPATH}/bin:/usr/local/go/bin:$PATH"

            if (branch == 'master') {
            container('sonar-scanner') {
              stage('Sonar version') {
                sh "/opt/sonar-scanner/bin/sonar-scanner --version"
                   }
                stage('Sonar scanner') {
                 sh "/opt/sonar-scanner/bin/sonar-scanner"
                  }
               }
            }

            stage('Go Dependencies') {
              if (branch == 'master' || branch == 'stable' || branch == 'va') {
                // sh "go get -v" Aggiungi quando hai moduli
              }
            }

            stage('Go Build')  {
              if (branch == 'master' || branch == 'stable' || branch == 'va') {
                sh "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o papagaio-api"
              }
            }

            stage('Creazione tar.gz') {
              if (branch == 'master' || branch == 'stable' || branch == 'va') {
                sh "mkdir dist && cp papagaio-api dist/"
                tarball = "papagaio-api-${version}.tar.gz"
                sh "tar -zcvf ${tarball} dist"
              }
            }

            stage('Deploy su Nexus') {
              if (branch == 'master' || branch == 'stable' || branch == 'va') {
                withCredentials([usernamePassword(credentialsId: 'nexus', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD')]) {
                  sh "curl -v -k -u ${USERNAME}:${PASSWORD} --upload-file ${tarball} https://nexus.sorintdev.it/repository/binaries/it.sorintdev.papagaio/papagaio-api-${version}.tar.gz"
                }
              }

              container('docker') {
                stage('Docker Image and Push') {
                  docker.withRegistry('https://registry.sorintdev.it', 'nexus') {
                    def img = docker.build(appName, '.')
                    if (branch == 'stable') {
                      img.push("${version}")
                    } else {
                      img.push("latest")
                    }
                  }
                if (branch == 'master' || branch == 'stable' || branch == 'va') {
                  stage ('Kubernetes') {
                      // sh "sed -i s/VERSION/${version}/g k8s/uat/deployment.yml"
                      // sh "sed -i s/VERSION/${version}/g k8s/va/deployment.yml"
                      sh "kubectl replace --force --kubeconfig=/etc/kubernetes/kubernetes.conf -f k8s/${target}"
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}