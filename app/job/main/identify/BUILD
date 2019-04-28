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
        "//app/job/main/identify/cmd:all-srcs",
        "//app/job/main/identify/conf:all-srcs",
        "//app/job/main/identify/dao:all-srcs",
        "//app/job/main/identify/http:all-srcs",
        "//app/job/main/identify/model:all-srcs",
        "//app/job/main/identify/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
