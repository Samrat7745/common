#!/bin/bash

#
# Copyright (c) 2021-2026 - for information on the respective copyright owner
# see the NOTICE file and/or the repository https://github.com/carbynestack/common.
#
# SPDX-License-Identifier: Apache-2.0
#

if [[ -z "${VERSION}" ]]; then
  echo "Artifact version must be defined in env variable \"VERSION\""
  exit -1
fi

if [[ -z  "${GITHUB_WORKSPACE}" ]]; then
  echo "Repository base directory must be defined in env variable \"GITHUB_WORKSPACE\""
  exit -2
fi

tmpChangelog=/tmp/cs.repository-changelog.${VERSION}
touch $tmpChangelog

appendGitLog() {
  recentTag=$(git describe --tags --abbrev=0)
  if (( $? != 0 )); then
    git log --all --format="- %s" --no-merges >> $1
  else
    git log ${recentTag}..HEAD --format="- %s" --no-merges >> $1
  fi
  echo -e "" >> $1
}

appendGoModuleRef() {
  if [ ! -f "${GITHUB_WORKSPACE}/go.mod" ]; then
    echo "Project is not a Go module"
    return
  fi

  module=$(cd ${GITHUB_WORKSPACE} && grep -m 1 -E "^module " go.mod | awk '{print $2}')
  if [[ -z "${module}" ]]; then
    echo "No Go module path found"
    return
  fi

  echo -e "## Go Module\n" >> $1
  echo -e "\`\`\`shell" >> $1
  echo -e "go get ${module}@${VERSION}" >> $1
  echo -e "\`\`\`" >> $1
}

echo -e "# Change Log\n" > $tmpChangelog

appendGitLog $tmpChangelog
appendGoModuleRef $tmpChangelog