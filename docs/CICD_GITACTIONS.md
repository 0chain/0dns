

## Guide to CI/CD using github actions
  <!-- Details of CI/CD setup using github -->
## Workflow Creation.
 - A new workflow is created using Go project with the file name called "build.yml".
 - By default the path of build.yml is ".github/workflows.build.yml"
 - Completed or running CI/CD can be seen under actions option.


## Details of components being used in build.yml.
#### Workflow name
Here the name of the workflow is defined i.e. "Dockerize"
```
name: Dockerize
```

#### Input Option to trigger manually builds
To run the workflow using manual option, *work_dispatch* is used. Which will ask for the input to tigger the builds with *latest* tag or not. If we select for **yes**, image will be build with *latest* tag as well as with *branch-commitid* tag. But if we select for **no**, image will be build with *branch-commitid* tag only.

```
on:
  workflow_dispatch:
    inputs:
      latest_tag:
        description: 'type yes for building latest tag'
        default: 'no'
        required: true
```

#### Global ENV setup
Environment variable is defined with the secrets added to the repository. Here secrets contains the docker images(example- dockerhub) repository name.
```
env:
  DNS_REGISTRY: ${{ secrets.DNS_REGISTRY }}
```

#### Defining jobs and runner
Jobs are defined which contains the various steps for creating and pushing the builds. Runner envionment is also defined used for making the builds.
```
jobs:
   dockerize_dns:
       runs-on: ubuntu-20.04
```

#### Different steps used in creating the builds
Here different steps are defined used for creating the builds.
 - *uses* --> checkout to branch from what code to create the builds.
 - *Get the version* --> Creating the tags by combining the branch name & first 8 digits of commit id.
 - *Login to Docker Hub* --> Logging into the docker hub using Username and Password from secrets of the repository.
 - *Build zdns* --> Building, tagging and pushing the docker images with the *Get the version* tag.
 - *Push image* --> Here we are checking if the input given by user is **yes**, images is also pushed with latest tag also.
```
- uses: actions/checkout@v2

- name: Get the version
    id: get_version
    run: |
    BRANCH=$(echo ${GITHUB_REF#refs/heads/} | sed 's/\//-/g')
    SHORT_SHA=$(echo $GITHUB_SHA | head -c 8)
    echo ::set-output name=BRANCH::${BRANCH}
    echo ::set-output name=VERSION::${BRANCH}-${SHORT_SHA}    
- name: Login to Docker Hub
    uses: docker/login-action@v1
    with:
    username: ${{ secrets.DOCKERHUB_USERNAME }}
    password: ${{ secrets.DOCKERHUB_PASSWORD }}

- name: Build zdns
    run: |
    docker build -t $DNS_REGISTRY:$TAG -f "$DOCKERFILE_DNS" .
    docker tag $DNS_REGISTRY:$TAG $DNS_REGISTRY:latest
    docker push $DNS_REGISTRY:latest
    env:
    TAG: ${{ steps.get_version.outputs.VERSION }}
    DOCKERFILE_DNS: "docker.local/Dockerfile"

- name: Push image 
    run: |
    if [[ "$PUSH_LATEST" == "yes" ]]; then
        docker push $DNS_REGISTRY:latest
    else
        docker push $DNS_REGISTRY:$TAG
    fi
    env:
    PUSH_LATEST: ${{ github.event.inputs.latest_tag }}
    TAG: ${{ steps.get_version.outputs.VERSION }}
```

