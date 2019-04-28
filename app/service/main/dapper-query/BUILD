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
        "//app/service/main/dapper-query/api/v1:all-srcs",
        "//app/service/main/dapper-query/cmd:all-srcs",
        "//app/service/main/dapper-query/conf:all-srcs",
        "//app/service/main/dapper-query/dao:all-srcs",
        "//app/service/main/dapper-query/model:all-srcs",
        "//app/service/main/dapper-query/pkg/cltclient:all-srcs",
        "//app/service/main/dapper-query/pkg/opslog:all-srcs",
        "//app/service/main/dapper-query/service:all-srcs",
        "//app/service/main/dapper-query/util:all-srcs",
    ],
    tags = ["automanaged"],
)
