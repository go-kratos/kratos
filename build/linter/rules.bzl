load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_repository")
load("//build/linter:deps.bzl", "DEPS")
load("//build/linter:internal/linter_repositories.bzl", "linter_repositories")

def go_lint_repositories(
    lang_deps = DEPS,
    lang_requires = [
        "com_github_alecthomas_gometalinter",
        "com_github_tsenart_deadcode",
        "com_github_mdempsky_maligned",
        "com_github_mibk_dupl",
        "com_github_kisielk_errcheck",
        "com_github_goastscanner_gas",
        "com_github_jgautheron_goconst",
        "com_github_alecthomas_gocyclo",
        "org_golang_x_goimports",
        "com_github_golang_lint",
        "co_honnef_tools_gosimple",
        "org_golang_x_gotype",
        "com_github_gordonklaus_ineffassign",
    ], **kwargs):

  rem = linter_repositories(lang_deps = lang_deps,
                           lang_requires = lang_requires,
                           **kwargs)

  # Load remaining (special) deps
  for dep in rem:
    rule = dep.pop("rule")
    if "go_repository" == rule:
      go_repository(**dep)
    else:
      fail("Unknown loading rule %s for %s" % (rule, dep))