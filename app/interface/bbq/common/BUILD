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
        "//app/interface/bbq/common/auth:all-srcs",
        "//app/interface/bbq/common/http:all-srcs",
        "//app/interface/bbq/common/model:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
