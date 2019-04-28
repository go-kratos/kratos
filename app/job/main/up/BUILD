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
        "//app/job/main/up/cmd:all-srcs",
        "//app/job/main/up/conf:all-srcs",
        "//app/job/main/up/dao/account:all-srcs",
        "//app/job/main/up/dao/email:all-srcs",
        "//app/job/main/up/dao/upcrm:all-srcs",
        "//app/job/main/up/model:all-srcs",
        "//app/job/main/up/server/http:all-srcs",
        "//app/job/main/up/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
