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
        "//app/interface/bbq/topic/cmd:all-srcs",
        "//app/interface/bbq/topic/internal/dao:all-srcs",
        "//app/interface/bbq/topic/internal/model:all-srcs",
        "//app/interface/bbq/topic/internal/server/http:all-srcs",
        "//app/interface/bbq/topic/internal/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
