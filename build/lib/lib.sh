set -o errexit
set -o nounset
set -o pipefail

# Unset CDPATH so that path interpolation can work correctly
# https://github.com/kratosrnetes/kratosrnetes/issues/52255
unset CDPATH

# The root of the build/dist directory
if [ -z "$KRATOS_ROOT" ]
then   
    KRATOS_ROOT="$(cd "$(dirname "${BASH_SOURCE}")/../.." && pwd -P)"
if