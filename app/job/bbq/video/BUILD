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
        "//app/job/bbq/video/cmd:all-srcs",
        "//app/job/bbq/video/conf:all-srcs",
        "//app/job/bbq/video/dao:all-srcs",
        "//app/job/bbq/video/model:all-srcs",
        "//app/job/bbq/video/server/http:all-srcs",
        "//app/job/bbq/video/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
