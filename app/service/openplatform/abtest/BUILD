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
        "//app/service/openplatform/abtest/cmd:all-srcs",
        "//app/service/openplatform/abtest/conf:all-srcs",
        "//app/service/openplatform/abtest/dao:all-srcs",
        "//app/service/openplatform/abtest/http:all-srcs",
        "//app/service/openplatform/abtest/model:all-srcs",
        "//app/service/openplatform/abtest/service:all-srcs",
    ],
    tags = ["automanaged"],
)
