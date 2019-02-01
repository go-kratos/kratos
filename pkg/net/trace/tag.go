package trace

// Standard Span tags https://github.com/opentracing/specification/blob/master/semantic_conventions.md#span-tags-table
const (
	// The software package, framework, library, or module that generated the associated Span.
	// E.g., "grpc", "django", "JDBI".
	// type string
	TagComponent = "component"

	// Database instance name.
	// E.g., In java, if the jdbc.url="jdbc:mysql://127.0.0.1:3306/customers", the instance name is "customers".
	// type string
	TagDBInstance = "db.instance"

	// A database statement for the given database type.
	// E.g., for db.type="sql", "SELECT * FROM wuser_table"; for db.type="redis", "SET mykey 'WuValue'".
	TagDBStatement = "db.statement"

	// Database type. For any SQL database, "sql". For others, the lower-case database category,
	// e.g. "cassandra", "hbase", or "redis".
	// type string
	TagDBType = "db.type"

	// Username for accessing database. E.g., "readonly_user" or "reporting_user"
	// type string
	TagDBUser = "db.user"

	// true if and only if the application considers the operation represented by the Span to have failed
	// type bool
	TagError = "error"

	// HTTP method of the request for the associated Span. E.g., "GET", "POST"
	// type string
	TagHTTPMethod = "http.method"

	// HTTP response status code for the associated Span. E.g., 200, 503, 404
	// type integer
	TagHTTPStatusCode = "http.status_code"

	// URL of the request being handled in this segment of the trace, in standard URI format.
	// E.g., "https://domain.net/path/to?resource=here"
	// type string
	TagHTTPURL = "http.url"

	// An address at which messages can be exchanged.
	// E.g. A Kafka record has an associated "topic name" that can be extracted by the instrumented producer or consumer and stored using this tag.
	// type string
	TagMessageBusDestination = "message_bus.destination"

	// Remote "address", suitable for use in a networking client library.
	// This may be a "ip:port", a bare "hostname", a FQDN, or even a JDBC substring like "mysql://prod-db:3306"
	// type string
	TagPeerAddress = "peer.address"

	// 	Remote hostname. E.g., "opentracing.io", "internal.dns.name"
	// type string
	TagPeerHostname = "peer.hostname"

	// Remote IPv4 address as a .-separated tuple. E.g., "127.0.0.1"
	// type string
	TagPeerIPv4 = "peer.ipv4"

	// Remote IPv6 address as a string of colon-separated 4-char hex tuples.
	// E.g., "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
	// type string
	TagPeerIPv6 = "peer.ipv6"

	// Remote port. E.g., 80
	// type integer
	TagPeerPort = "peer.port"

	// Remote service name (for some unspecified definition of "service").
	// E.g., "elasticsearch", "a_custom_microservice", "memcache"
	// type string
	TagPeerService = "peer.service"

	// If greater than 0, a hint to the Tracer to do its best to capture the trace.
	// If 0, a hint to the trace to not-capture the trace. If absent, the Tracer should use its default sampling mechanism.
	// type string
	TagSamplingPriority = "sampling.priority"

	// Either "client" or "server" for the appropriate roles in an RPC,
	// and "producer" or "consumer" for the appropriate roles in a messaging scenario.
	// type string
	TagSpanKind = "span.kind"

	// legacy tag
	TagAnnotation = "legacy.annotation"
	TagAddress    = "legacy.address"
	TagComment    = "legacy.comment"
)

//  Standard log tags
const (
	// The type or "kind" of an error (only for event="error" logs). E.g., "Exception", "OSError"
	// type string
	LogErrorKind = "error.kind"

	// For languages that support such a thing (e.g., Java, Python),
	// the actual Throwable/Exception/Error object instance itself.
	// E.g., A java.lang.UnsupportedOperationException instance, a python exceptions.NameError instance
	// type string
	LogErrorObject = "error.object"

	// A stable identifier for some notable moment in the lifetime of a Span. For instance, a mutex lock acquisition or release or the sorts of lifetime events in a browser page load described in the Performance.timing specification. E.g., from Zipkin, "cs", "sr", "ss", or "cr". Or, more generally, "initialized" or "timed out". For errors, "error"
	// type string
	LogEvent = "event"

	// A concise, human-readable, one-line message explaining the event.
	// E.g., "Could not connect to backend", "Cache invalidation succeeded"
	// type string
	LogMessage = "message"

	// A stack trace in platform-conventional format; may or may not pertain to an error. E.g., "File \"example.py\", line 7, in \<module\>\ncaller()\nFile \"example.py\", line 5, in caller\ncallee()\nFile \"example.py\", line 2, in callee\nraise Exception(\"Yikes\")\n"
	// type string
	LogStack = "stack"
)

// Tag interface
type Tag struct {
	Key   string
	Value interface{}
}

// TagString new string tag.
func TagString(key string, val string) Tag {
	return Tag{Key: key, Value: val}
}

// TagInt64 new int64 tag.
func TagInt64(key string, val int64) Tag {
	return Tag{Key: key, Value: val}
}

// TagInt new int tag
func TagInt(key string, val int) Tag {
	return Tag{Key: key, Value: val}
}

// TagBool new bool tag
func TagBool(key string, val bool) Tag {
	return Tag{Key: key, Value: val}
}

// TagFloat64 new float64 tag
func TagFloat64(key string, val float64) Tag {
	return Tag{Key: key, Value: val}
}

// TagFloat32 new float64 tag
func TagFloat32(key string, val float32) Tag {
	return Tag{Key: key, Value: val}
}

// String new tag String.
// NOTE: use TagString
func String(key string, val string) Tag {
	return TagString(key, val)
}

// Int new tag Int.
// NOTE: use TagInt
func Int(key string, val int) Tag {
	return TagInt(key, val)
}

// Bool new tagBool
// NOTE: use TagBool
func Bool(key string, val bool) Tag {
	return TagBool(key, val)
}

// Log new log.
func Log(key string, val string) LogField {
	return LogField{Key: key, Value: val}
}

// LogField LogField
type LogField struct {
	Key   string
	Value string
}
