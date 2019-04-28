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
        "//app/job/openplatform/article/cmd:all-srcs",
        "//app/job/openplatform/article/conf:all-srcs",
        "//app/job/openplatform/article/dao:all-srcs",
        "//app/job/openplatform/article/http:all-srcs",
        "//app/job/openplatform/article/model:all-srcs",
        "//app/job/openplatform/article/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
