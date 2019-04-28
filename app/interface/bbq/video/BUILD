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
        "//app/interface/bbq/video/cmd:all-srcs",
        "//app/interface/bbq/video/internal/dao:all-srcs",
        "//app/interface/bbq/video/internal/model:all-srcs",
        "//app/interface/bbq/video/internal/server/http:all-srcs",
        "//app/interface/bbq/video/internal/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
