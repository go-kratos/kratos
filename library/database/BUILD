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
        "//library/database/bfs:all-srcs",
        "//library/database/elastic:all-srcs",
        "//library/database/hbase.v2:all-srcs",
        "//library/database/orm:all-srcs",
        "//library/database/sql:all-srcs",
        "//library/database/tidb:all-srcs",
    ],
    tags = ["automanaged"],
)
