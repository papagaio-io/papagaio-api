local appName = "papagaio-api";
local targetMap = {master: "dev", stable: "stable", release: "release"};
local versionNumber = "0.1.1";

local go_runtime() = {
  type: 'pod',
  arch: 'amd64',
  containers: [
    { image: 'registry.sorintdev.it/golang:1.15' },
  ],
};

local task_build_go(branch) = 
  local version = if branch == "master" then "latest" else versionNumber;
  local tarball = appName + "-" + version + ".tar.gz";
  {
    name: "build go " + branch,
    runtime: go_runtime(),
    when: {
      branch: {
        include: if branch == "" then [] else [branch],
        exclude: if branch == "" then ["master", "stable", "release"] else [],
      }
    },
    environment:
    {
      "USERNAME": {
        from_variable: "NEXUSUSERNAME"
      },
       "PASSWORD": {
        from_variable: "NEXUSPASSWORD"
      },
    },
    steps: 
      [
        { type: 'clone' },
        { type: 'restore_cache', keys: ['cache-sum-{{ md5sum "go.sum" }}', 'cache-date-'], dest_dir: '/go/pkg/mod/cache' },
        { type: 'run', name: 'build the program', command: 'go build .' },
        { type: 'save_cache', key: 'cache-sum-{{ md5sum "go.sum" }}', contents: [{ source_dir: '/go/pkg/mod/cache' }] },
        { type: 'save_cache', key: 'cache-date-{{ year }}-{{ month }}-{{ day }}', contents: [{ source_dir: '/go/pkg/mod/cache' }] },
        { type: 'save_to_workspace', contents: [{ source_dir: '.', dest_dir: '.', paths: ['**'] }] },
      ] +
      if branch == "master" || branch == "stable" then
      [
        { type: 'run', name: 'Creazione tar.gz', command: 'mkdir dist && cp papagaio-api dist/ && tar -zcvf ' + tarball + ' dist' },
        { type: 'run', name: 'Deploy su Nexus', command: 'curl -v -k -u $USERNAME:$PASSWORD --upload-file ' + tarball + ' https://nexus.sorintdev.it/repository/binaries/it.sorintdev.papagaio/papagaio-api-' + version + '.tar.gz' },
      ] else [],
  };

local task_docker_build_push(branch) = {
  name: 'docker build and push ' + branch,
  when: {
    branch: [branch],
  },
  runtime:
  {
    containers: [
      { image: "gcr.io/kaniko-project/executor:debug"},
    ],
  },
  environment:
  {
    "DOCKERAUTH": {
      from_variable: "dockerauth"
    },
    "VERSION_NUMBER": versionNumber,
    "APPNAME": appName,
  },
  shell: "/busybox/sh",
  working_dir: '/kaniko',
  steps: 
   [
    { type: "restore_workspace", name: "restore workspace", dest_dir: "/kaniko/papagaio-api" },
    { type: 'run', name: 'test1', command: 'ls -la /kaniko/papagaio-api' },
    {
      type: "run",
      name: "generate docker config", 
      command: |||
        cat << EOF > /kaniko/.docker/config.json
        {
          "auths": {
            "registry.sorintdev.it": { "auth": "$DOCKERAUTH" }
          }
        }
        EOF
      |||,
    },
    {
      type: "run",
      name: "kanico executor",
      command: |||
        echo "branch" $AGOLA_GIT_BRANCH
        if [[ $AGOLA_GIT_BRANCH == 'master' ]]; then 
          export VERSION="latest" ;
        else
          export VERSION=${VERSION_NUMBER} ; fi
        echo "version" $VERSION
        /kaniko/executor --context=dir:///kaniko/papagaio-api --dockerfile Dockerfile --destination registry.sorintdev.it/$APPNAME:$VERSION
      |||,
    },
   ],
  depends: ["build go " + branch]
};

local task_kubernetes_deploy(branch) = 
  local version = if branch == "master" then "latest" else versionNumber;
  local target = targetMap[branch];
  {
    name: "kubernetes deploy " + branch,
    runtime:
    {
      containers: [
        { 
          image: "registry.sorintdev.it/bitnami/kubectl:1.19",
          volumes: [
            {
              path: "/mnt/data",
              tmpfs: {},
            },
          ],
        },
      ],
    },
    environment:
    {
      "KUBERNETESCONF": {
        from_variable: "SORINTDEVKUBERNETESCONF"
      },
    },
    when: {
      branch: branch
    },
    working_dir: '/mnt/data',
    steps: 
    [
      { type: "restore_workspace", name: "restore workspace", dest_dir: "." },
      { type: 'run', name: 'create folder kubernetes', command: 'mkdir kubernetes' },
      { type: 'run', name: 'generate kubernetes config', command: 'echo $KUBERNETESCONF | base64 -d > kubernetes/kubernetes.conf' },
      { type: 'run', name: 'sed version stable', command: 'sed -i s/VERSION/' + version + '/g k8s/stable/deployment.yml' },
      { type: 'run', name: 'sed version release', command: 'sed -i s/VERSION/' + version + '/g k8s/release/deployment.yml' },
      { type: 'run', name: 'kubectl replace', command: 'kubectl replace --force --kubeconfig=kubernetes/kubernetes.conf -f k8s/' + target },
    ],
    depends: if branch == "release" then ["build go release"] else ["docker build and push " + branch],
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
            from_variable: "NEXUSUSERNAME"
          },
          password:
          {
            from_variable: "NEXUSPASSWORD"
          }
        }
      },
      tasks: [
        task_build_go("master"),
        task_build_go("stable"),
        task_build_go("release"),
        task_build_go(""),
        task_docker_build_push("master"),
        task_docker_build_push("stable"),
        task_kubernetes_deploy("master"),
        task_kubernetes_deploy("stable"),
        task_kubernetes_deploy("release"),
      ]
    },
  ],
}
