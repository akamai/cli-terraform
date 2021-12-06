#!/usr/bin/env bash

clean() {
  if [[ -d "bin" ]]; then
    rm -rf bin
  fi;

  if [[ -d "test" ]]; then
    rm -rf test
  fi;

}

find_branch() {
  CURRENT_BRANCH=$GIT_BRANCH
  if [[ "$CURRENT_BRANCH" == "develop" ]]; then
    EDGEGRID_BRANCH="v2"
  elif [[ $CURRENT_BRANCH =~ .*/sp-.* ]]; then
    EDGEGRID_BRANCH=${CURRENT_BRANCH//origin\//}
  else
    # find parent branch from which this branch was created, iterate over the list of branches from the history of commits
    branches=($(git log --pretty=format:'%D' | sed 's@HEAD -> @@' | grep . | sed 's@origin/@@g' | sed -e $'s@, @\\\n@g' | grep -v HEAD ))
    for branch in ${branches[*]}
    do
      echo "Checking branch '${branch}'"
      EDGEGRID_BRANCH=$branch
      if [[ "$EDGEGRID_BRANCH" == "develop" ]]; then
        EDGEGRID_BRANCH="v2"
      fi
      git -C ./akamaiopen-edgegrid-golang branch -r | grep $EDGEGRID_BRANCH > /dev/null
      if [[ $? -eq 0 ]]; then
        echo "There is matching EdgeGrid branch '${EDGEGRID_BRANCH}'"
        break
      fi
    done
  fi
  echo "Current branch is '${CURRENT_BRANCH}', matching EdgeGrid branch is '${EDGEGRID_BRANCH}'"
}

clone_edgegrid() {
  if [ ! -d "./akamaiopen-edgegrid-golang" ]
  then
    echo "First time build, cloning the 'akamaiopen-edgegrid-golang' repo"
    git clone ssh://git@git.source.akamai.com:7999/devexp/akamaiopen-edgegrid-golang.git
  else
    echo "Repository 'akamaiopen-edgegrid-golang' already exists, so only cleaning and updating it"
    pushd akamaiopen-edgegrid-golang
    git reset --hard
    git fetch --prune
    popd
  fi
}

checkout_edgegrid() {
  pushd akamaiopen-edgegrid-golang
  git checkout $EDGEGRID_BRANCH -f
  git reset --hard origin/$EDGEGRID_BRANCH
  git pull
  popd
}

clean
clone_edgegrid
find_branch

if ! ./build/docker_jenkins.bash "$CURRENT_BRANCH" "$CURRENT_BRANCH" "$EDGEGRID_BRANCH"; then
    exit 1
fi
