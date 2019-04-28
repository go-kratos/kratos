package tomllint

var synataxerrordata = `
[owner]
name = "Tom Preston-Werner"
dob 1979-05-27T07:32:00-08:00 # First class dates

`
var normaldata = `
[owner]
name = "Tom Preston-Werner"
dob = 1979-05-27T07:32:00-08:00 # First class dates

[database]
server = "192.168.1.1"
ports = [ 8001, 8001, 8002 ]
connection_max = 5000
enabled = true

[servers]

  # Indentation (tabs and/or spaces) is allowed but not required
  [servers.alpha]
  ip = "10.0.0.1"
  dc = "eqdc10"

  [servers.beta]
  ip = "10.0.0.2"
  dc = "eqdc10"

[clients]
data = [ ["gamma", "delta"], [1, 2] ]

# Line breaks are OK when inside arrays
hosts = [
  "alpha",
  "omega"
]
`

var notopkvdata = `
title = "test123"
age = 123
array = [1,2,3]
[clients]
data = "helo"
`

var nocommondata = `
version = "v1.1"
user = "nobody"
`

var noidentify = `
[identify]
xxx = 123
`

var noapp = `
[app]
xxx = 123
`
