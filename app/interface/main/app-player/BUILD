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
        "//app/interface/main/app-player/cmd:all-srcs",
        "//app/interface/main/app-player/conf:all-srcs",
        "//app/interface/main/app-player/dao:all-srcs",
        "//app/interface/main/app-player/http:all-srcs",
        "//app/interface/main/app-player/model:all-srcs",
        "//app/interface/main/app-player/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
