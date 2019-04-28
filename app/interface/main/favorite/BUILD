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
        "//app/interface/main/favorite/cmd:all-srcs",
        "//app/interface/main/favorite/conf:all-srcs",
        "//app/interface/main/favorite/dao/music:all-srcs",
        "//app/interface/main/favorite/dao/topic:all-srcs",
        "//app/interface/main/favorite/dao/video:all-srcs",
        "//app/interface/main/favorite/http:all-srcs",
        "//app/interface/main/favorite/model:all-srcs",
        "//app/interface/main/favorite/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
