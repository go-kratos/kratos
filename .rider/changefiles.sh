#!/bin/bash

#set -x

suffix=${1}

url="${CI_SERVER_URL}/api/v4/projects/${CI_PROJECT_ID}/repository/commits/${CI_COMMIT_SHA}/merge_requests?private_token=${CI_PRIVATE_TOKEN}"
json=$(curl -s ${url})
if [[ "$json" != "[]" && "$json" != "" ]]; then
    target_branch=$(echo ${json} | jq -r '.[0].target_branch')
    source_branch=$(echo ${json} | jq -r '.[0].source_branch')
    files=$(git diff origin/${target_branch}...origin/${source_branch} --name-only  --diff-filter=ACM | grep -E -i "${suffix}")
else
#  url="${CI_SERVER_URL}/api/v4/projects/${CI_PROJECT_ID}/pipelines?private_token=${CI_PRIVATE_TOKEN}&status=success&ref=${CI_COMMIT_REF_NAME}"
#  commit=$(curl -s ${url} | jq -r 'first(.[] | .sha)')
##  echo "Last green commit is '${commit}'."
#  files=$(git diff ${commit} --name-only  --diff-filter=ACM | grep -E "${suffix}")
    files=$(git diff origin/master...origin/${CI_COMMIT_REF_NAME} --name-only  --diff-filter=ACM | grep -E -i "${suffix}")
fi

echo -e "${files}"