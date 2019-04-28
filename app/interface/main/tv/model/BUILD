load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "audit.go",
        "audit_v2.go",
        "auth.go",
        "favorite.go",
        "follow.go",
        "loadpage.go",
        "media.go",
        "media_v2.go",
        "module.go",
        "others.go",
        "pgc_cards.go",
        "playurl.go",
        "recommend.go",
        "region.go",
        "upgrade.go",
        "zone_index.go",
    ],
    importpath = "go-common/app/interface/main/tv/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/archive/api:go_default_library",
        "//library/ecode:go_default_library",
        "//library/time:go_default_library",
    ],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [
        ":package-srcs",
        "//app/interface/main/tv/model/goblin:all-srcs",
        "//app/interface/main/tv/model/history:all-srcs",
        "//app/interface/main/tv/model/search:all-srcs",
        "//app/interface/main/tv/model/thirdp:all-srcs",
        "//app/interface/main/tv/model/tvvip:all-srcs",
        "//app/interface/main/tv/model/upper:all-srcs",
        "//app/interface/main/tv/model/view:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
