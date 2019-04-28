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
        "//app/job/main/app-player/cmd:all-srcs",
        "//app/job/main/app-player/conf:all-srcs",
        "//app/job/main/app-player/dao:all-srcs",
        "//app/job/main/app-player/http:all-srcs",
        "//app/job/main/app-player/model:all-srcs",
        "//app/job/main/app-player/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
