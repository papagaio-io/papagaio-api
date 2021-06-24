local appName = "papagaio-api"
local targetMap = [ master: 'dev', stable: 'stable', release: 'release']
local branch = AGOLA_GIT_BRANCH
local target = targetMap[branch]
local label = appName + "-${UUID.randomUUID().toString()}"
local version = "0.1.1"

if branch == master then
   version = "latest"
;

local go_runtime() = {
  type: 'pod',
  arch: 'amd64',
  containers: [
    { image: 'registry.sorintdev.it/golang:1.15.8' },
  ],
};

tarball = "papagaio-api-" + version + ".tar.gz"

local task_build_go() = {
  name: 'build go',
  runtime: go_runtime(),
  steps: [
    { type: 'clone' },
    { type: 'restore_cache', keys: ['cache-sum-{{ md5sum "go.sum" }}', 'cache-date-'], dest_dir: '/go/pkg/mod/cache' },
    { type: 'run', name: 'build the program', command: 'go build .' },
    { type: 'save_to_workspace', contents: [{ source_dir: '.', dest_dir: '/bin/', paths: ['agola-example-go'] }] },
    { type: 'save_cache', key: 'cache-sum-{{ md5sum "go.sum" }}', contents: [{ source_dir: '/go/pkg/mod/cache' }] },
    { type: 'save_cache', key: 'cache-date-{{ year }}-{{ month }}-{{ day }}', contents: [{ source_dir: '/go/pkg/mod/cache' }] },
    //TODO CREAZIONE tar.gz (TODO aggiungere if branch==master/stable/release)
    { type: 'run', name: 'Creazione tar.gz', command: 'mkdir dist && cp papagaio-api dist/ && tar -zcvf ' + tarball + ' dist' },
    //TODO deploy tar.gz su nexus (TODO aggiungere if branch==master/stable)
    { type: 'run', name: 'Deploy su Nexus', command: '"curl -v -k -u ${SORINT_DOCKER_USERNAME}:${SORINT_DOCKER_PASSWORD} --upload-file ' + tarball + " https://nexus.sorintdev.it/repository/binaries/it.sorintdev.papagaio/papagaio-api-" + version + ".tar.gz" }
  ],
};

local task_docker_build_push() = {
  name: 'docker build and push',
  runtime:
  {
    containers:
    {
      image: "gcr.io/kaniko-project/executor:debug"
    }
  },
  environment:
  {
    DOCKERAUTH: {
      from_variable: "dockerauth"
    }
  },
  shell: "/busybox/sh"
  steps: 
  {
    restore_workspace: 
    {
      dest_dir: "."
    },
    run: "/kaniko/executor --destination registry/image"
  },
  depends: "build go"
};

{
  runs: [
    {
      name: 'papagaio backend',
      docker_registries_auth:
      {
        'registry.sorintdev.it':
        {
          username:
          {
            from_variable: "SORINT_DOCKER_USERNAME" //TODO agola variable from secret
          },
          password:
          {
            from_variable: "SORINT_DOCKER_PASSWORD" //TODO agola variable from secret
          }
        }
      },
      tasks: [
        task_build_go(),
        task_docker_build_push()
      ]
    },
  ],
}
