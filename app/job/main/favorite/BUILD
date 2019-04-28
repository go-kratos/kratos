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
        "//app/job/main/favorite/cmd:all-srcs",
        "//app/job/main/favorite/conf:all-srcs",
        "//app/job/main/favorite/dao/fav:all-srcs",
        "//app/job/main/favorite/dao/music:all-srcs",
        "//app/job/main/favorite/dao/pub:all-srcs",
        "//app/job/main/favorite/dao/stat:all-srcs",
        "//app/job/main/favorite/http:all-srcs",
        "//app/job/main/favorite/model:all-srcs",
        "//app/job/main/favorite/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
