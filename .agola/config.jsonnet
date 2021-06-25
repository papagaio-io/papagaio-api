local appName = "papagaio-api";
local targetMap = {master: "dev", stable: "stable", release: "release"};
local branch = std.extVar("AGOLA_GIT_BRANCH");
local target = targetMap[branch];
local label = appName + "-${UUID.randomUUID().toString()}";

local versionNumber = "0.1.1";
local versionMap = {master: "laster", "stable": versionNumber, "release": versionNumber};
local version = versionMap[branch];

local go_runtime() = {
  type: 'pod',
  arch: 'amd64',
  containers: [
    { image: 'registry.sorintdev.it/golang:1.15.8' },
  ],
};

local tarball = "papagaio-api-" + version + ".tar.gz";

local task_build_go() = {
  name: 'build go',
  runtime: go_runtime(),
  steps: 
    [
      { type: 'clone' },
      { type: 'restore_cache', keys: ['cache-sum-{{ md5sum "go.sum" }}', 'cache-date-'], dest_dir: '/go/pkg/mod/cache' },
      { type: 'run', name: 'build the program', command: 'go build .' },
      { type: 'save_to_workspace', contents: [{ source_dir: '.', dest_dir: '/bin/', paths: ['papagaio-api'] }] },
      { type: 'save_cache', key: 'cache-sum-{{ md5sum "go.sum" }}', contents: [{ source_dir: '/go/pkg/mod/cache' }] },
      { type: 'save_cache', key: 'cache-date-{{ year }}-{{ month }}-{{ day }}', contents: [{ source_dir: '/go/pkg/mod/cache' }] }
    ] +
    if branch == "master" || branch == "stable" || branch == "release" then
    [
      { type: 'run', name: 'Creazione tar.gz', command: 'mkdir dist && cp papagaio-api dist/ && tar -zcvf ' + tarball + ' dist' },
    ] +
    if branch == master || branch == "stable" then
    [
      { type: 'run', name: 'Deploy su Nexus', command: '"curl -v -k -u $sorint-docker-username:$sorint-docker-password --upload-file ' + tarball + " https://nexus.sorintdev.it/repository/binaries/it.sorintdev.papagaio/papagaio-api-" + version + ".tar.gz" }
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
    DOCKERAUTHUSERNAME: {
      from_variable: "sorint-docker-username"
    },
    DOCKERAUTHPASSWORD: {
      from_variable: "dockersorint-docker-password"
    },
  },
  shell: "/busybox/sh",
  steps: 
  {
    restore_workspace: 
    {
      dest_dir: "."
    },
    run: [
      {
        name: "generate docker config",
        command: |||
          cat << EOF > /kaniko/.docker/config.json
                {
                  "auths": {
                    "registry.sorintdev.it": { "auth" : "$DOCKERAUTHUSERNAME:$DOCKERAUTHPASSWORD" }
                  }
                }
          EOF
          |||,
      },
        "/kaniko/executor --destination registry.sorintdev.it/" + appName + ":" + version
    ]
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
            from_variable: "sorint-docker-username"
          },
          password:
          {
            from_variable: "sorint-docker-password"
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
