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
        "//app/service/openplatform/ticket-item/api/grpc/v1:all-srcs",
        "//app/service/openplatform/ticket-item/cmd:all-srcs",
        "//app/service/openplatform/ticket-item/conf:all-srcs",
        "//app/service/openplatform/ticket-item/dao:all-srcs",
        "//app/service/openplatform/ticket-item/model:all-srcs",
        "//app/service/openplatform/ticket-item/server/grpc:all-srcs",
        "//app/service/openplatform/ticket-item/server/http:all-srcs",
        "//app/service/openplatform/ticket-item/service:all-srcs",
    ],
    tags = ["automanaged"],
)
