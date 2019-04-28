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
        "//app/interface/main/player/cmd:all-srcs",
        "//app/interface/main/player/conf:all-srcs",
        "//app/interface/main/player/dao:all-srcs",
        "//app/interface/main/player/http:all-srcs",
        "//app/interface/main/player/model:all-srcs",
        "//app/interface/main/player/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
