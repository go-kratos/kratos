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
        "//app/service/main/sms/api:all-srcs",
        "//app/service/main/sms/cmd:all-srcs",
        "//app/service/main/sms/conf:all-srcs",
        "//app/service/main/sms/dao:all-srcs",
        "//app/service/main/sms/model:all-srcs",
        "//app/service/main/sms/server/gorpc:all-srcs",
        "//app/service/main/sms/server/grpc:all-srcs",
        "//app/service/main/sms/server/http:all-srcs",
        "//app/service/main/sms/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
