local appName = "papagaio";

local go_runtime() = {
  type: 'pod',
  arch: 'amd64',
  containers: [
    { image: 'registry.sorintdev.it/golang:1.15' },
  ],
};

local task_test() = 
  {
    name: "test",
    runtime: go_runtime(),
    steps: 
      [
        { type: 'clone' },
        { type: 'restore_cache', keys: ['cache-sum-{{ md5sum "go.sum" }}', 'cache-date-'], dest_dir: '/go/pkg/mod/cache' },
        { type: 'run', name: 'go unit tests', command: 'go test -coverprofile testCover.out ./service' },
        { type: 'save_cache', key: 'cache-sum-{{ md5sum "go.sum" }}', contents: [{ source_dir: '/go/pkg/mod/cache' }] },
        { type: 'save_cache', key: 'cache-date-{{ year }}-{{ month }}-{{ day }}', contents: [{ source_dir: '/go/pkg/mod/cache' }] },
      ],
  };

local task_build_go() = 
  {
    name: "build go",
    runtime: go_runtime(),
    depends: ["test"],
    environment:
    {
      "USERNAME": {
        from_variable: "NEXUS-USERNAME"
      },
       "PASSWORD": {
        from_variable: "NEXUS-PASSWORD"
      },
        "url_repo_upload": {
        from_variable: "URL-REPO-UPLOAD"
      },
    },
    steps: 
      [
        { type: 'clone' },
        { type: 'save_to_workspace', contents: [{ source_dir: '.', dest_dir: '.', paths: ['**'] }] },
        { type: 'restore_cache', keys: ['cache-sum-{{ md5sum "go.sum" }}', 'cache-date-'], dest_dir: '/go/pkg/mod/cache' },
        { type: 'run', command: 'make' },
        { type: 'save_cache', key: 'cache-sum-{{ md5sum "go.sum" }}', contents: [{ source_dir: '/go/pkg/mod/cache' }] },
        { type: 'save_cache', key: 'cache-date-{{ year }}-{{ month }}-{{ day }}', contents: [{ source_dir: '/go/pkg/mod/cache' }] },
        { type: 'save_to_workspace', contents: [{ source_dir: './bin', dest_dir: '/bin/', paths: ['*'] }] },
        { type: 'run',
          name: 'Create and deploy Nexus', 
          command: |||
            export
            if [ ${AGOLA_GIT_TAG} ]; then
              export TARBALL=papagaio-${AGOLA_GIT_TAG}.tar.gz ;
            else
              export TARBALL=papagaio-latest.tar.gz ; fi

            mkdir dist && cp bin/papagaio dist/ && tar -zcvf ${TARBALL} dist
            curl -v -k -u $USERNAME:$PASSWORD --upload-file ${TARBALL} ${url_repo_upload}${TARBALL}
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
    arch: 'amd64',
    containers: [
      { image: "gcr.io/kaniko-project/executor:debug-v0.11.0"},
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
    { type: "restore_workspace", name: "restore workspace", dest_dir: "/kaniko/papagaio" },
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
          /kaniko/executor --context=dir:///kaniko/papagaio --build-arg PAPAGAIOWEB_IMAGE=tulliobotti/papagaio-web:v2.0.2 --target papagaio --dockerfile Dockerfile --destination registry.sorintdev.it/$APPNAME:$AGOLA_GIT_TAG;
        else
          /kaniko/executor --context=dir:///kaniko/papagaio --build-arg PAPAGAIOWEB_IMAGE=tulliobotti/papagaio-web:v2.0.2 --target papagaio --dockerfile Dockerfile --destination registry.sorintdev.it/$APPNAME:latest ; fi
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
    { type: "restore_workspace", name: "restore workspace", dest_dir: "/kaniko/papagaio" },
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
    { type: "run", name: "kanico executor", command: "/kaniko/executor --context=dir:///kaniko/papagaio --build-arg PAPAGAIOWEB_IMAGE=tulliobotti/papagaio-web:v2.0.2 --dockerfile Dockerfile --destination tulliobotti/$APPNAME:$AGOLA_GIT_TAG" },
   ],
  depends: ["build go"]
};

local task_kubernetes_deploy(namespace, changeImageVersion) = 
  {
    name: "kubernetes deploy " + namespace,
    environment:
    {
      "KUBERNETESCONF": {
        from_variable: "SORINT-DEV-KUBERNETES-CONF"
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
        { type: 'run', name: 'kubectl replace', command:  'kubectl --kubeconfig=kubernetes/kubernetes.conf -n ' + namespace + ' set image deployment/papagaio papagaio=registry.sorintdev.it/papagaio:$AGOLA_GIT_TAG' }
      else 
        { type: 'run', name: 'kubectl replace', command:  'kubectl --kubeconfig=kubernetes/kubernetes.conf -n ' + namespace + ' delete pods -l app=papagaio' },
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
            from_variable: "NEXUS-USERNAME"
          },
          password:
          {
            from_variable: "NEXUS-PASSWORD"
          }
        }
      },
      tasks: [
        task_test(),
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