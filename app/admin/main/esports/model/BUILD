package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "archive.go",
        "contest.go",
        "contest_data.go",
        "game.go",
        "gid_map.go",
        "match.go",
        "match_active.go",
        "match_detail.go",
        "match_map.go",
        "season.go",
        "tag.go",
        "tag_map.go",
        "team.go",
        "team_map.go",
        "tree.go",
        "year_map.go",
    ],
    importpath = "go-common/app/admin/main/esports/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = ["//library/xstr:go_default_library"],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
