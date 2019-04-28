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
        "//app/job/main/upload/cmd:all-srcs",
        "//app/job/main/upload/conf:all-srcs",
        "//app/job/main/upload/dao:all-srcs",
        "//app/job/main/upload/http:all-srcs",
        "//app/job/main/upload/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
