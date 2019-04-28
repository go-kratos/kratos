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
        "//app/common/openplatform/geetest/conf:all-srcs",
        "//app/common/openplatform/geetest/dao:all-srcs",
        "//app/common/openplatform/geetest/model:all-srcs",
        "//app/common/openplatform/geetest/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
