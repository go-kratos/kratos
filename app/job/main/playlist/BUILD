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
        "//app/job/main/playlist/cmd:all-srcs",
        "//app/job/main/playlist/conf:all-srcs",
        "//app/job/main/playlist/dao:all-srcs",
        "//app/job/main/playlist/http:all-srcs",
        "//app/job/main/playlist/model:all-srcs",
        "//app/job/main/playlist/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
