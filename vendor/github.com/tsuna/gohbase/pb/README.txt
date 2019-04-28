These are the protobuf definition files used by GoHBase.
They were copied from HBase (see under hbase-protocol/src/main/protobuf).

Currently using .proto files from HBase branch-1.3, with unnecessary .proto
files removed.

The following changes were made to those files:
  - the package name was changed to "pb". (`sed -i 's/hbase.pb/pb/g' ./*`)

The files in this directory are also subject to the Apache License 2.0 and
are copyright of the Apache Software Foundation.
