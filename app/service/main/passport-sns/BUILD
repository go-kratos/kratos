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
        "//app/service/main/passport-sns/api:all-srcs",
        "//app/service/main/passport-sns/cmd:all-srcs",
        "//app/service/main/passport-sns/conf:all-srcs",
        "//app/service/main/passport-sns/dao:all-srcs",
        "//app/service/main/passport-sns/model:all-srcs",
        "//app/service/main/passport-sns/server/http:all-srcs",
        "//app/service/main/passport-sns/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
