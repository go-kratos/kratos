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
        "//app/service/main/riot-search/benchmark:all-srcs",
        "//app/service/main/riot-search/cmd:all-srcs",
        "//app/service/main/riot-search/conf:all-srcs",
        "//app/service/main/riot-search/dao:all-srcs",
        "//app/service/main/riot-search/model:all-srcs",
        "//app/service/main/riot-search/server/http:all-srcs",
        "//app/service/main/riot-search/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
