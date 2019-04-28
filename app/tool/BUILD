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
        "//app/tool/bgr:all-srcs",
        "//app/tool/bmproto/bmgen:all-srcs",
        "//app/tool/bmproto/protoc-gen-bm:all-srcs",
        "//app/tool/cache:all-srcs",
        "//app/tool/ci:all-srcs",
        "//app/tool/creator:all-srcs",
        "//app/tool/gdoc:all-srcs",
        "//app/tool/gengo/args:all-srcs",
        "//app/tool/gengo/cmd/deepcopy-gen:all-srcs",
        "//app/tool/gengo/generator:all-srcs",
        "//app/tool/gengo/namer:all-srcs",
        "//app/tool/gengo/parser:all-srcs",
        "//app/tool/gengo/testdata/a:all-srcs",
        "//app/tool/gengo/testdata/fake/dep:all-srcs",
        "//app/tool/gengo/types:all-srcs",
        "//app/tool/gorpc:all-srcs",
        "//app/tool/grpc-http-proxy:all-srcs",
        "//app/tool/kratos:all-srcs",
        "//app/tool/liverpc/liverpcgen:all-srcs",
        "//app/tool/liverpc/protoc-gen-liverpc:all-srcs",
        "//app/tool/mkprow:all-srcs",
        "//app/tool/owner:all-srcs",
        "//app/tool/protoc-gen-bm:all-srcs",
        "//app/tool/saga:all-srcs",
        "//app/tool/warden:all-srcs",
    ],
    tags = ["automanaged"],
)
