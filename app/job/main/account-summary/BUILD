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
        "//app/job/main/account-summary/cmd:all-srcs",
        "//app/job/main/account-summary/conf:all-srcs",
        "//app/job/main/account-summary/dao:all-srcs",
        "//app/job/main/account-summary/http:all-srcs",
        "//app/job/main/account-summary/model:all-srcs",
        "//app/job/main/account-summary/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
