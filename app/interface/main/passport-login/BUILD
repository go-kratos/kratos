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
        "//app/interface/main/passport-login/cmd:all-srcs",
        "//app/interface/main/passport-login/conf:all-srcs",
        "//app/interface/main/passport-login/dao:all-srcs",
        "//app/interface/main/passport-login/model:all-srcs",
        "//app/interface/main/passport-login/server/http:all-srcs",
        "//app/interface/main/passport-login/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
