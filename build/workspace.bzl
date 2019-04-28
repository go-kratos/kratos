# TensorFlow external dependencies that can be loaded in WORKSPACE files.

load("//vendor:repo.bzl", "bili_http_archive")


# Sanitize a dependency so that it works correctly from code that includes
# TensorFlow as a submodule.
def clean_dep(dep):
  return str(Label(dep))

# If TensorFlow is linked as a submodule.
# path_prefix is no longer used.
# tf_repo_name is thought to be under consideration.
def bili_workspace(path_prefix="", tf_repo_name=""):
  # Note that we check the minimum bazel version in WORKSPACE.
  bili_http_archive(
      name = "pcre",
      sha256 = "84c3c4d2eb9166aaed44e39b89e4b6a49eac6fed273bdb844c94fb6c8bdda1b5",
      urls = [
          "http://bazel-cabin.bilibili.co/clib/pcre-8.42.zip",
      ],
      strip_prefix = "pcre-8.42",
      build_file = clean_dep("//vendor:pcre.BUILD"),
  )
