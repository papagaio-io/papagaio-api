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
            if [ ${AGOLA_GIT_TAG} ]; then
              export TARBALL=papagaio-api-${AGOLA_GIT_TAG}.tar.gz ;
            else
              export TARBALL=papagaio-api-latest.tar.gz ; fi

            mkdir dist && cp papagaio-api dist/ && tar -zcvf ${TARBALL} dist
            curl -v -k -u $USERNAME:$PASSWORD --upload-file ${TARBALL} https://nexus.sorintdev.it/repository/binaries/it.sorintdev.papagaio/${TARBALL}
          |||,
        },
      ],
  };

local task_docker_build_push() = {
  name: 'docker build and push',
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
    "PUBLIC_DOCKERAUTH": {
      from_variable: "SORINTLAB_DOCKERAUTH"
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
            "registry.sorintdev.it": { "auth": "$PRIVATE_DOCKERAUTH" },
            "sorintlab": { "auth": "$PUBLIC_DOCKERAUTH" }
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
        if [ ${AGOLA_GIT_TAG} ]; then
          /kaniko/executor --context=dir:///kaniko/papagaio-api --dockerfile Dockerfile --destination registry.sorintdev.it/$APPNAME:$AGOLA_GIT_TAG;
          /kaniko/executor --context=dir:///kaniko/papagaio-api --dockerfile Dockerfile --destination sorintlab/$APPNAME:$AGOLA_GIT_TAG;
        else
          /kaniko/executor --context=dir:///kaniko/papagaio-api --dockerfile Dockerfile --destination registry.sorintdev.it/$APPNAME:latest ; fi
      |||,
    },
   ],
  depends: ["build go"]
};

local task_kubernetes_deploy(target) = 
  {
    name: "kubernetes deploy " + target,
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
    working_dir: '/mnt/data',
    steps: 
    [
      { type: "restore_workspace", name: "restore workspace", dest_dir: "." },
      { type: 'run', name: 'create folder kubernetes', command: 'mkdir kubernetes' },
      { type: 'run', name: 'generate kubernetes config', command: 'echo $KUBERNETESCONF | base64 -d > kubernetes/kubernetes.conf' },
      { type: 'run', name: 'sed version stable', command: 'sed -i s/VERSION/$AGOLA_GIT_TAG/g k8s/stable/deployment.yml' },
      { type: 'run', name: 'sed version release', command: 'sed -i s/VERSION/$AGOLA_GIT_TAG/g k8s/release/deployment.yml' },
      { type: 'run', name: 'kubectl replace', command: 'kubectl replace --force --kubeconfig=kubernetes/kubernetes.conf -f k8s/' + target },
    ],
    depends: ["docker build and push"],
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
        task_docker_build_push(),
        task_kubernetes_deploy('dev')+ {
          when: {
            branch: 'master',
          },
        },
        task_kubernetes_deploy('stable')+ {
          when: {
            tag: '#.*#',
          },
        },
        task_kubernetes_deploy('release')+ {
          when: {
            tag: '#.*#',
          },
          approval: true,
        },
      ]
    },
  ],
}
