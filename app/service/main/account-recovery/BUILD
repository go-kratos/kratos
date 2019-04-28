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
        "//app/service/main/account-recovery/cmd:all-srcs",
        "//app/service/main/account-recovery/conf:all-srcs",
        "//app/service/main/account-recovery/dao:all-srcs",
        "//app/service/main/account-recovery/http:all-srcs",
        "//app/service/main/account-recovery/model:all-srcs",
        "//app/service/main/account-recovery/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
