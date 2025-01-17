#!/usr/bin/env bash

clean() {
  if [[ -d "bin" ]]; then
    rm -rf bin
  fi;

  if [[ -d "test" ]]; then
    rm -rf test
  fi;
}

list_branches() {
  git log --pretty=format:'%D' | sed 's@HEAD -> @@' | grep . | sed 's@origin/@@g' | sed 's@release/.*@@g' | sed -E $'s@master, (.+)@\\1, master@g' | tr ', ' '\n' | grep -v 'tag:' | sed -E 's@^v([0-9]+\.?){2,}(-rc\.[0-9]+)?@@g' | grep -v release/ | grep -v HEAD | sed '/^$/d'
}

find_edgegrid_branch() {
  CURRENT_BRANCH=$GIT_BRANCH
  if [[ "$CURRENT_BRANCH" == "develop" || "$CURRENT_BRANCH" == "master" ]]; then
    echo Current branch is '${CURRENT_BRANCH}'
    EDGEGRID_BRANCH=${CURRENT_BRANCH//origin\//}
  else
    # find parent branch from which this branch was created, iterate over the list of branches from the history of commits
    branches=($(list_branches))
    branches+=("develop") # guard to fallback to safe value if less branches than 5
    for branch in ${branches[*]}
    do
      echo "Checking branch '${branch}'"
      EDGEGRID_BRANCH=$branch

      if [[ "$index" -eq "5" ]]; then
        echo "Exceeding limit of checks, fallback to default branch 'develop'"
        EDGEGRID_BRANCH="develop"
        break
      fi
      index=$((index + 1))

      git -C ./akamaiopen-edgegrid-golang branch -r | grep $EDGEGRID_BRANCH > /dev/null
      if [[ $? -eq 0 ]]; then
        echo "There is matching EdgeGrid branch '${EDGEGRID_BRANCH}'"
        break
      fi
    done
  fi
  echo "Current branch is '${CURRENT_BRANCH}', matching EdgeGrid branch is '${EDGEGRID_BRANCH}'"
}

find_cli_branch() {
  CURRENT_BRANCH=$GIT_BRANCH
  # find parent branch from which this branch was created, iterate over the list of branches from the history of commits
  branches=($(list_branches))
  branches+=("develop") # guard to fallback to safe value if less branches than 5
  for branch in ${branches[*]}
  do
    echo "Checking Cli branch '${branch}'"
    CLI_BRANCH=$branch

    if [[ "$index" -eq "5" ]]; then
      echo "Exceeding limit of checks, fallback to default branch 'develop'"
      CLI_BRANCH="develop"
      break
    fi
    index=$((index + 1))

    git -C ./cli-clone branch -r | grep $CLI_BRANCH > /dev/null
    if [[ $? -eq 0 ]]; then
      echo "There is matching Cli branch '${CLI_BRANCH}'"
      break
    fi
  done
  echo "Current branch is '${CURRENT_BRANCH}', matching Cli branch is '${CLI_BRANCH}'"
}

clone_repository() {
  case "$1" in
    edgegrid)
      repo="akamaiopen-edgegrid-golang"
      ;;
    provider)
      repo="terraform-provider-akamai"
      ;;
    cli)
      repo="cli"
      ;;
    *)
      echo "Repository '${1}' is unknown, exiting..." && exit 1
      ;;
  esac
  target_dir=${2:-$repo}

  if [ ! -d "./${target_dir}" ]
  then
    echo "First time build, cloning the '${repo}' repo into '${target_dir}'"
    git clone ssh://git@git.source.akamai.com:7999/devexp/${repo}.git $target_dir
  else
    echo "Repository '${repo}' already exists, so only cleaning and updating it"
    pushd ${target_dir}
    git reset --hard
    git fetch --prune
    popd
  fi
}

clean
clone_repository edgegrid
clone_repository cli cli-clone
find_edgegrid_branch
find_cli_branch

if ! ./build/docker_jenkins.bash "$CURRENT_BRANCH" "$EDGEGRID_BRANCH" "$CLI_BRANCH"; then
    exit 1
fi
