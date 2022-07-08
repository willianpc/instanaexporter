package instanaexporter

// AttributeInstanaHostID can be used to distinguish multiple hosts' data
// being processed by a single collector (in a chained scenario)
const AttributeInstanaHostID = "instana.host.id"

const HeaderKey = "x-instana-key"
const HeaderHost = "x-instana-host"
const HeaderTime = "x-instana-time"
