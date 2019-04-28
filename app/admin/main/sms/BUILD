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
        "//app/admin/main/sms/cmd:all-srcs",
        "//app/admin/main/sms/conf:all-srcs",
        "//app/admin/main/sms/dao:all-srcs",
        "//app/admin/main/sms/model:all-srcs",
        "//app/admin/main/sms/server/http:all-srcs",
        "//app/admin/main/sms/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
