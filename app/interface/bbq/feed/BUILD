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
        "//app/interface/bbq/feed/cmd:all-srcs",
        "//app/interface/bbq/feed/internal/dao:all-srcs",
        "//app/interface/bbq/feed/internal/model:all-srcs",
        "//app/interface/bbq/feed/internal/server/http:all-srcs",
        "//app/interface/bbq/feed/internal/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
