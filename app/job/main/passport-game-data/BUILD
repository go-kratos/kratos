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
        "//app/job/main/passport-game-data/cmd:all-srcs",
        "//app/job/main/passport-game-data/conf:all-srcs",
        "//app/job/main/passport-game-data/dao:all-srcs",
        "//app/job/main/passport-game-data/http:all-srcs",
        "//app/job/main/passport-game-data/model:all-srcs",
        "//app/job/main/passport-game-data/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
