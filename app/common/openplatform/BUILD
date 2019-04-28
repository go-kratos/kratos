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
        "//app/common/openplatform/encoding:all-srcs",
        "//app/common/openplatform/geetest:all-srcs",
        "//app/common/openplatform/random:all-srcs",
    ],
    tags = ["automanaged"],
)
