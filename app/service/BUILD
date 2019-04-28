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
        "//app/service/bbq:all-srcs",
        "//app/service/ep:all-srcs",
        "//app/service/live:all-srcs",
        "//app/service/main:all-srcs",
        "//app/service/openplatform:all-srcs",
        "//app/service/ops/log-agent:all-srcs",
        "//app/service/video/stream-mng/api/v1:all-srcs",
        "//app/service/video/stream-mng/cmd:all-srcs",
        "//app/service/video/stream-mng/common:all-srcs",
        "//app/service/video/stream-mng/conf:all-srcs",
        "//app/service/video/stream-mng/dao:all-srcs",
        "//app/service/video/stream-mng/middleware:all-srcs",
        "//app/service/video/stream-mng/model:all-srcs",
        "//app/service/video/stream-mng/server/grpc:all-srcs",
        "//app/service/video/stream-mng/server/http:all-srcs",
        "//app/service/video/stream-mng/service:all-srcs",
    ],
    tags = ["automanaged"],
)
