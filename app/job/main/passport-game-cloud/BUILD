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
        "//app/job/main/passport-game-cloud/cmd:all-srcs",
        "//app/job/main/passport-game-cloud/conf:all-srcs",
        "//app/job/main/passport-game-cloud/dao:all-srcs",
        "//app/job/main/passport-game-cloud/http:all-srcs",
        "//app/job/main/passport-game-cloud/model:all-srcs",
        "//app/job/main/passport-game-cloud/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
