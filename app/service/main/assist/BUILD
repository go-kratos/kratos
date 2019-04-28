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
        "//app/service/main/assist/cmd:all-srcs",
        "//app/service/main/assist/conf:all-srcs",
        "//app/service/main/assist/dao/account:all-srcs",
        "//app/service/main/assist/dao/assist:all-srcs",
        "//app/service/main/assist/dao/message:all-srcs",
        "//app/service/main/assist/http:all-srcs",
        "//app/service/main/assist/model:all-srcs",
        "//app/service/main/assist/rpc:all-srcs",
        "//app/service/main/assist/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
