package(default_visibility = ["//visibility:public"])

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
        "//app/service/main/passport-game/cmd:all-srcs",
        "//app/service/main/passport-game/conf:all-srcs",
        "//app/service/main/passport-game/dao:all-srcs",
        "//app/service/main/passport-game/http:all-srcs",
        "//app/service/main/passport-game/model:all-srcs",
        "//app/service/main/passport-game/service:all-srcs",
    ],
    tags = ["automanaged"],
)
