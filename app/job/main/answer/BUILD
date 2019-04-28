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
        "//app/job/main/answer/cmd:all-srcs",
        "//app/job/main/answer/conf:all-srcs",
        "//app/job/main/answer/dao:all-srcs",
        "//app/job/main/answer/http:all-srcs",
        "//app/job/main/answer/model:all-srcs",
        "//app/job/main/answer/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
