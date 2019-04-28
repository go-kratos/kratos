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
        "//app/service/openplatform/anti-fraud/api/grpc/v1:all-srcs",
        "//app/service/openplatform/anti-fraud/api/http/v1:all-srcs",
        "//app/service/openplatform/anti-fraud/cmd:all-srcs",
        "//app/service/openplatform/anti-fraud/conf:all-srcs",
        "//app/service/openplatform/anti-fraud/dao:all-srcs",
        "//app/service/openplatform/anti-fraud/model:all-srcs",
        "//app/service/openplatform/anti-fraud/server/grpc:all-srcs",
        "//app/service/openplatform/anti-fraud/server/http:all-srcs",
        "//app/service/openplatform/anti-fraud/service:all-srcs",
    ],
    tags = ["automanaged"],
)
