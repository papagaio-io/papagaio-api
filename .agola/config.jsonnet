local appName = "papagaio-api";

local go_runtime() = {
  type: 'pod',
  arch: 'amd64',
  containers: [
    { image: 'registry.sorintdev.it/golang:1.15' },
  ],
};

local task_build_go() = 
  {
    name: "build go",
    runtime: go_runtime(),
    environment:
    {
      "USERNAME": {
        from_variable: "NEXUSUSERNAME"
      },
       "PASSWORD": {
        from_variable: "NEXUSPASSWORD"
      },
        "urlrepoupload": {
        from_variable: "URLREPOUPLOAD"
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
        { type: 'run',
          name: 'Create and deploy Nexus', 
          command: |||
            export
            if [ ${AGOLA_GIT_TAG} ]; then
              export TARBALL=papagaio-api-${AGOLA_GIT_TAG}.tar.gz ;
            else
              export TARBALL=papagaio-api-latest.tar.gz ; fi

            mkdir dist && cp papagaio-api dist/ && tar -zcvf ${TARBALL} dist
            curl -v -k -u $USERNAME:$PASSWORD --upload-file ${TARBALL} ${urlrepoupload}${TARBALL}
          |||,
        },
      ],
  };

local task_docker_build_push_private() = {
  name: 'docker build and push private',
  when: {
    branch: 'master',
    tag: '#.*#',
  },
  runtime:
  {
    containers: [
      { image: "gcr.io/kaniko-project/executor:debug"},
    ],
  },
  environment:
  {
    "PRIVATE_DOCKERAUTH": {
      from_variable: "dockerauth"
    },
    "APPNAME": appName,
  },
  shell: "/busybox/sh",
  working_dir: '/kaniko',
  steps: 
   [
    { type: "restore_workspace", name: "restore workspace", dest_dir: "/kaniko/papagaio-api" },
    {
      type: "run",
      name: "generate docker config", 
      command: |||
        cat << EOF > /kaniko/.docker/config.json
        {
          "auths": {
            "registry.sorintdev.it": { "auth": "$PRIVATE_DOCKERAUTH" }
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
        if [ $AGOLA_GIT_TAG ]; then
          /kaniko/executor --context=dir:///kaniko/papagaio-api --dockerfile Dockerfile --destination registry.sorintdev.it/$APPNAME:$AGOLA_GIT_TAG;
        else
          /kaniko/executor --context=dir:///kaniko/papagaio-api --dockerfile Dockerfile --destination registry.sorintdev.it/$APPNAME:latest ; fi
      |||,
    },
   ],
  depends: ["build go"]
};

local task_docker_build_push_public() = {
  name: 'docker build and push public',
  when: {
    tag: '#.*#',
  },
  runtime:
  {
    containers: [
      { image: "gcr.io/kaniko-project/executor:debug"},
    ],
  },
  environment:
  {
    "APPNAME": appName,
    "PUBLIC_DOCKERAUTH": {
      from_variable: "TULLIO-DOCKERAUTH"
    },
  },
  shell: "/busybox/sh",
  working_dir: '/kaniko',
  steps: 
   [
    { type: "restore_workspace", name: "restore workspace", dest_dir: "/kaniko/papagaio-api" },
    {
      type: "run",
      name: "generate docker config", 
      command: |||
        cat << EOF > /kaniko/.docker/config.json
        {
          "auths": {
            "https://index.docker.io/v1/": { "auth": "$PUBLIC_DOCKERAUTH" }
          }
        }
        EOF
      |||,
    },
    { type: "run", name: "kanico executor", command: "/kaniko/executor --context=dir:///kaniko/papagaio-api --dockerfile Dockerfile --destination tulliobotti/$APPNAME:$AGOLA_GIT_TAG" },
   ],
  depends: ["build go"]
};

local task_kubernetes_deploy(namespace, changeImageVersion) = 
  {
    name: "kubernetes deploy " + namespace,
    environment:
    {
      "KUBERNETESCONF": {
        from_variable: "SORINTDEVKUBERNETESCONF"
      },
    },
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
    working_dir: '/mnt/data',
    steps:
    [
      { type: "restore_workspace", name: "restore workspace", dest_dir: "." },
      { type: 'run', name: 'create folder kubernetes', command: 'mkdir kubernetes' },
      { type: 'run', name: 'generate kubernetes config', command: 'echo $KUBERNETESCONF | base64 -d > kubernetes/kubernetes.conf' },
            
      if changeImageVersion then
        { type: 'run', name: 'kubectl replace', command:  'kubectl --kubeconfig=kubernetes/kubernetes.conf -n ' + namespace + ' set image deployment/papagaio-api papagaio-api=registry.sorintdev.it/papagaio-api:$AGOLA_GIT_TAG' }
      else 
        { type: 'run', name: 'kubectl replace', command:  'kubectl --kubeconfig=kubernetes/kubernetes.conf -n ' + namespace + ' delete pods -l app=papagaio-api' },
    ],
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
        task_build_go(),
        task_docker_build_push_private(),
        task_docker_build_push_public(),
        task_kubernetes_deploy('ci-dev', false)+ {
          when: {
            branch: 'master',
          },
          depends: ["docker build and push private"],
        },
        task_kubernetes_deploy('ci-stable', true)+ {
          when: {
            tag: '#.*#',
          },
          depends: ["docker build and push private", "docker build and push public"],
        },
        task_kubernetes_deploy('ci', true)+ {
          when: {
            tag: '#.*#',
          },
          approval: true,
          depends: ["docker build and push private", "docker build and push public"],
        },
      ]
    },
  ],
}