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
        "//app/interface/openplatform/seo/cmd:all-srcs",
        "//app/interface/openplatform/seo/conf:all-srcs",
        "//app/interface/openplatform/seo/dao:all-srcs",
        "//app/interface/openplatform/seo/model:all-srcs",
        "//app/interface/openplatform/seo/server/http:all-srcs",
        "//app/interface/openplatform/seo/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
