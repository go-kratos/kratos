package(default_visibility = ["//visibility:public"])

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
        "//app/service/main/secure/cmd:all-srcs",
        "//app/service/main/secure/conf:all-srcs",
        "//app/service/main/secure/dao:all-srcs",
        "//app/service/main/secure/http:all-srcs",
        "//app/service/main/secure/model:all-srcs",
        "//app/service/main/secure/rpc/client:all-srcs",
        "//app/service/main/secure/rpc/server:all-srcs",
        "//app/service/main/secure/service:all-srcs",
    ],
    tags = ["automanaged"],
)
