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
        "//app/interface/main/app-tag/cmd:all-srcs",
        "//app/interface/main/app-tag/conf:all-srcs",
        "//app/interface/main/app-tag/dao/archive:all-srcs",
        "//app/interface/main/app-tag/dao/bangumi:all-srcs",
        "//app/interface/main/app-tag/dao/recommend:all-srcs",
        "//app/interface/main/app-tag/dao/region:all-srcs",
        "//app/interface/main/app-tag/dao/tag:all-srcs",
        "//app/interface/main/app-tag/http:all-srcs",
        "//app/interface/main/app-tag/model:all-srcs",
        "//app/interface/main/app-tag/service/ping:all-srcs",
        "//app/interface/main/app-tag/service/tag:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
