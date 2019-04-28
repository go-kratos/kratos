load("//build/linter:internal/require.bzl", "require")
load("//build/linter:deps.bzl", "DEPS")

def linter_repositories(excludes = [],
                       lang_deps = {},
                       lang_requires = [],
                       overrides = {},
                       strict = False,
                       verbose = 0):
  return require(
    keys =  lang_requires,
    deps =  lang_deps,
    excludes = excludes,
    overrides = overrides,
    verbose = verbose,
    strict = strict,
)