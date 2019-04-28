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
        "//app/job/openplatform/article:all-srcs",
        "//app/job/openplatform/open-market:all-srcs",
        "//app/job/openplatform/open-sug/api/grpc/v1:all-srcs",
        "//app/job/openplatform/open-sug/cmd:all-srcs",
        "//app/job/openplatform/open-sug/conf:all-srcs",
        "//app/job/openplatform/open-sug/dao:all-srcs",
        "//app/job/openplatform/open-sug/model:all-srcs",
        "//app/job/openplatform/open-sug/server/http:all-srcs",
        "//app/job/openplatform/open-sug/service:all-srcs",
    ],
    tags = ["automanaged"],
)
