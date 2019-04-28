load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "audit.go",
        "auth.go",
        "history.go",
        "http.go",
        "media.go",
        "others.go",
        "pages.go",
        "recommend.go",
        "search.go",
        "tvvip.go",
        "ugc.go",
        "ugc_view.go",
        "upgrade.go",
        "wild.go",
        "zone_index.go",
    ],
    importpath = "go-common/app/interface/main/tv/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/tv/conf:go_default_library",
        "//app/interface/main/tv/dao/thirdp:go_default_library",
        "//app/interface/main/tv/model:go_default_library",
        "//app/interface/main/tv/model/search:go_default_library",
        "//app/interface/main/tv/model/tvvip:go_default_library",
        "//app/interface/main/tv/service/app:go_default_library",
        "//app/interface/main/tv/service/audit:go_default_library",
        "//app/interface/main/tv/service/favorite:go_default_library",
        "//app/interface/main/tv/service/goblin:go_default_library",
        "//app/interface/main/tv/service/history:go_default_library",
        "//app/interface/main/tv/service/pgc:go_default_library",
        "//app/interface/main/tv/service/search:go_default_library",
        "//app/interface/main/tv/service/thirdp:go_default_library",
        "//app/interface/main/tv/service/tvvip:go_default_library",
        "//app/interface/main/tv/service/view:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/binding:go_default_library",
        "//library/net/http/blademaster/middleware/auth:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
        "//library/net/http/blademaster/render:go_default_library",
        "//library/net/metadata:go_default_library",
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
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
