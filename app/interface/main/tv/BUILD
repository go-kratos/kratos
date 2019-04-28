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
        "//app/interface/main/tv/cmd:all-srcs",
        "//app/interface/main/tv/conf:all-srcs",
        "//app/interface/main/tv/dao/account:all-srcs",
        "//app/interface/main/tv/dao/app:all-srcs",
        "//app/interface/main/tv/dao/archive:all-srcs",
        "//app/interface/main/tv/dao/audit:all-srcs",
        "//app/interface/main/tv/dao/cms:all-srcs",
        "//app/interface/main/tv/dao/favorite:all-srcs",
        "//app/interface/main/tv/dao/goblin:all-srcs",
        "//app/interface/main/tv/dao/history:all-srcs",
        "//app/interface/main/tv/dao/pgc:all-srcs",
        "//app/interface/main/tv/dao/search:all-srcs",
        "//app/interface/main/tv/dao/thirdp:all-srcs",
        "//app/interface/main/tv/dao/upper:all-srcs",
        "//app/interface/main/tv/http:all-srcs",
        "//app/interface/main/tv/model:all-srcs",
        "//app/interface/main/tv/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
