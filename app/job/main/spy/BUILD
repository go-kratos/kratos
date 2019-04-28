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
        "//app/job/main/spy/cmd:all-srcs",
        "//app/job/main/spy/conf:all-srcs",
        "//app/job/main/spy/dao:all-srcs",
        "//app/job/main/spy/http:all-srcs",
        "//app/job/main/spy/model:all-srcs",
        "//app/job/main/spy/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
