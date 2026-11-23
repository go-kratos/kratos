#!/usr/bin/env bash

# Create release tags for every contrib Go module in this repository.
# A module's tag prefix follows its physical repository directory.

set -o errexit
set -o nounset
set -o pipefail

usage() {
    echo "Usage: $0 VERSION [--push]" >&2
    echo "Example: $0 v3.0.0 --push" >&2
}

if [[ "$#" -lt 1 || "$#" -gt 2 ]]; then
    usage
    exit 2
fi

version="$1"
push_tags=false
if [[ "$#" -eq 2 ]]; then
    if [[ "$2" != "--push" ]]; then
        usage
        exit 2
    fi
    push_tags=true
fi

if [[ ! "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+([.-][0-9A-Za-z.-]+)?$ ]]; then
    echo "invalid release version: ${version}" >&2
    exit 1
fi

KRATOS_HOME=$(
    cd "$(dirname "${BASH_SOURCE[0]}")" &&
    cd .. &&
    pwd
)
cd "${KRATOS_HOME}"

release_commit=$(git rev-parse --verify HEAD^{commit})
module_major="${version%%.*}"
module_files=()
while IFS= read -r -d '' module_file; do
    module_files+=("${module_file}")
done < <(find contrib -type f -name go.mod -print0 | sort -z)

if [[ "${#module_files[@]}" -eq 0 ]]; then
    echo "no contrib Go modules found" >&2
    exit 1
fi

release_tags=()
for module_file in "${module_files[@]}"; do
    module_dir="${module_file%/go.mod}"
    module_path=$(awk '/^module[[:space:]]/ { module = $2; sub(/\r$/, "", module); print module; exit }' "${module_file}")

    if [[ "${module_path}" != */"${module_major}" ]]; then
        continue
    fi

    tag="${module_dir}/${version}"
    release_tags+=("${tag}")
done

if [[ "${#release_tags[@]}" -eq 0 ]]; then
    echo "no contrib modules match release major ${module_major}" >&2
    exit 1
fi

for tag in "${release_tags[@]}"; do
    if git rev-parse --verify --quiet "refs/tags/${tag}^{commit}" >/dev/null; then
        existing_commit=$(git rev-parse "refs/tags/${tag}^{commit}")
        if [[ "${existing_commit}" != "${release_commit}" ]]; then
            echo "tag ${tag} already points to ${existing_commit}, expected ${release_commit}" >&2
            exit 1
        fi
    fi
done

for module_file in "${module_files[@]}"; do
    module_dir="${module_file%/go.mod}"
    module_path=$(awk '/^module[[:space:]]/ { module = $2; sub(/\r$/, "", module); print module; exit }' "${module_file}")

    if [[ "${module_path}" != */"${module_major}" ]]; then
        continue
    fi

    tag="${module_dir}/${version}"
    if git rev-parse --verify --quiet "refs/tags/${tag}^{commit}" >/dev/null; then
        echo "tag ${tag} already points to ${release_commit}"
        continue
    fi

    git tag --annotate "${tag}" "${release_commit}" --message "release ${module_path} ${version}"
    echo "created ${tag} -> ${release_commit}"
done

if [[ "${push_tags}" == true ]]; then
    git push origin "${release_tags[@]}"
fi
