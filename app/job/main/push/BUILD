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
        "//app/job/main/push/cmd:all-srcs",
        "//app/job/main/push/conf:all-srcs",
        "//app/job/main/push/dao:all-srcs",
        "//app/job/main/push/http:all-srcs",
        "//app/job/main/push/model:all-srcs",
        "//app/job/main/push/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
