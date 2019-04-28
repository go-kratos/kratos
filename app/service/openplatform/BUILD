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
        "//app/service/openplatform/abtest:all-srcs",
        "//app/service/openplatform/anti-fraud:all-srcs",
        "//app/service/openplatform/pgc-season:all-srcs",
        "//app/service/openplatform/ticket-item:all-srcs",
        "//app/service/openplatform/ticket-sales:all-srcs",
    ],
    tags = ["automanaged"],
)
