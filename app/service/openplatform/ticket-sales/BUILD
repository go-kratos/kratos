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
        "//app/service/openplatform/ticket-sales/api/grpc/type:all-srcs",
        "//app/service/openplatform/ticket-sales/api/grpc/v1:all-srcs",
        "//app/service/openplatform/ticket-sales/cmd:all-srcs",
        "//app/service/openplatform/ticket-sales/conf:all-srcs",
        "//app/service/openplatform/ticket-sales/dao:all-srcs",
        "//app/service/openplatform/ticket-sales/model:all-srcs",
        "//app/service/openplatform/ticket-sales/server/grpc:all-srcs",
        "//app/service/openplatform/ticket-sales/server/http:all-srcs",
        "//app/service/openplatform/ticket-sales/service:all-srcs",
    ],
    tags = ["automanaged"],
)
