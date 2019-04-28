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
        "//app/service/openplatform/pgc-season/api/grpc/episode:all-srcs",
        "//app/service/openplatform/pgc-season/api/grpc/season:all-srcs",
    ],
    tags = ["automanaged"],
)
