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
        "//app/interface/main/tv/service/app:all-srcs",
        "//app/interface/main/tv/service/audit:all-srcs",
        "//app/interface/main/tv/service/favorite:all-srcs",
        "//app/interface/main/tv/service/goblin:all-srcs",
        "//app/interface/main/tv/service/history:all-srcs",
        "//app/interface/main/tv/service/pgc:all-srcs",
        "//app/interface/main/tv/service/search:all-srcs",
        "//app/interface/main/tv/service/thirdp:all-srcs",
        "//app/interface/main/tv/service/tvvip:all-srcs",
        "//app/interface/main/tv/service/view:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
