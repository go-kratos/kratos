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
        "//app/interface/main/shorturl/cmd:all-srcs",
        "//app/interface/main/shorturl/conf:all-srcs",
        "//app/interface/main/shorturl/dao:all-srcs",
        "//app/interface/main/shorturl/http:all-srcs",
        "//app/interface/main/shorturl/model:all-srcs",
        "//app/interface/main/shorturl/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
