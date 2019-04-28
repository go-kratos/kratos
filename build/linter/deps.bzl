# ****************************************************************
# List of external dependencies
# ****************************************************************

# These deps are derived empirically, starting from the grpc tag and
# then updating other dependents to the master commit.

DEPS = {

    "com_github_alecthomas_gometalinter": {
        "rule": "go_repository",
        "importpath": "github.com/alecthomas/gometalinter",
        "tag": "v2.0.5",
    },

    "com_github_tsenart_deadcode": {
        "rule": "go_repository",
        "importpath": "github.com/tsenart/deadcode",
    },

    "com_github_mdempsky_maligned": {
        "rule": "go_repository",
        "importpath": "github.com/mdempsky/maligned",
    },

    "com_github_mibk_dupl": {
        "rule": "go_repository",
        "importpath": "github.com/mibk/dupl",
    },

    "com_github_kisielk_errcheck": {
        "rule": "go_repository",
        "importpath": "github.com/kisielk/errcheck",
    },

    "com_github_goastscanner_gas": {
        "rule": "go_repository",
        "importpath": "github.com/GoASTScanner/gas",
    },

    "com_github_jgautheron_goconst": {
        "rule": "go_repository",
        "importpath": "github.com/jgautheron/goconst/cmd/goconst",
    },

    "com_github_alecthomas_gocyclo": {
        "rule": "go_repository",
        "importpath": "github.com/alecthomas/gocyclo",
    },
    
    "org_golang_x_goimports": {
        "rule": "go_repository",
        "importpath": "golang.org/x/tools/cmd/goimports",
    },

    "com_github_golang_lint": {
        "rule": "go_repository",
        "importpath": "github.com/golang/lint/golint",
    },

    "co_honnef_tools_gosimple": {
        "rule": "go_repository",
        "importpath": "honnef.co/go/tools/cmd/gosimple",
    },

    "org_golang_x_gotype": {
        "rule": "go_repository",
        "importpath": "golang.org/x/tools/cmd/gotype",
    },

     "com_github_gordonklaus_ineffassign": {
        "rule": "go_repository",
        "importpath": "github.com/gordonklaus/ineffassign",
    },

}
