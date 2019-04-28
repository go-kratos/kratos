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
        "//app/service/main/share/api:all-srcs",
        "//app/service/main/share/cmd:all-srcs",
        "//app/service/main/share/conf:all-srcs",
        "//app/service/main/share/dao:all-srcs",
        "//app/service/main/share/model:all-srcs",
        "//app/service/main/share/server/grpc:all-srcs",
        "//app/service/main/share/server/http:all-srcs",
        "//app/service/main/share/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
