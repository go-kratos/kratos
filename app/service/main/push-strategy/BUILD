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
        "//app/service/main/push-strategy/cmd:all-srcs",
        "//app/service/main/push-strategy/conf:all-srcs",
        "//app/service/main/push-strategy/dao:all-srcs",
        "//app/service/main/push-strategy/http:all-srcs",
        "//app/service/main/push-strategy/model:all-srcs",
        "//app/service/main/push-strategy/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
