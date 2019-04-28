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
        "//app/job/main/sms/cmd:all-srcs",
        "//app/job/main/sms/conf:all-srcs",
        "//app/job/main/sms/dao:all-srcs",
        "//app/job/main/sms/http:all-srcs",
        "//app/job/main/sms/model:all-srcs",
        "//app/job/main/sms/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
