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
        "//app/job/main/dm/cmd:all-srcs",
        "//app/job/main/dm/conf:all-srcs",
        "//app/job/main/dm/dao:all-srcs",
        "//app/job/main/dm/http:all-srcs",
        "//app/job/main/dm/model:all-srcs",
        "//app/job/main/dm/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
