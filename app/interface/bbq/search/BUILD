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
        "//app/interface/bbq/search/cmd:all-srcs",
        "//app/interface/bbq/search/internal/dao:all-srcs",
        "//app/interface/bbq/search/internal/model:all-srcs",
        "//app/interface/bbq/search/internal/server/http:all-srcs",
        "//app/interface/bbq/search/internal/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
