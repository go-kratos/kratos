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
        "//app/job/main/passport-game-local/cmd:all-srcs",
        "//app/job/main/passport-game-local/conf:all-srcs",
        "//app/job/main/passport-game-local/http:all-srcs",
        "//app/job/main/passport-game-local/model:all-srcs",
        "//app/job/main/passport-game-local/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
