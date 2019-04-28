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
        "//app/job/main/account-recovery/cmd:all-srcs",
        "//app/job/main/account-recovery/conf:all-srcs",
        "//app/job/main/account-recovery/dao:all-srcs",
        "//app/job/main/account-recovery/model:all-srcs",
        "//app/job/main/account-recovery/server/http:all-srcs",
        "//app/job/main/account-recovery/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
