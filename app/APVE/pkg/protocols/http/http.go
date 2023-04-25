package http

type Request struct {
	Path []string `yaml:"path,omitempty" jsonschema:"title=path(s) for the http request,description=Path(s) to send http requests to"`
	Raw  []string `yaml:"raw,omitempty" jsonschema:"http requests in raw format,description=HTTP Requests in Raw Format"`
	// ID is the optional id of the request
	ID   string `yaml:"id,omitempty" jsonschema:"title=id for the http request,description=ID for the HTTP Request"`
	Name string `yaml:"name,omitempty" jsonschema:"title=name for the http request,description=Optional name for the HTTP Request"`
	Body string `yaml:"body,omitempty" jsonschema:"title=body is the http request body,description=Body is an optional parameter which contains HTTP Request body"`
	// description: |
	//   Payloads contains any payloads for the current request.
	//
	//   Payloads support both key-values combinations where a list
	//   of payloads is provided, or optionally a single file can also
	//   be provided as payload which will be read on run-time.
	Payloads map[string]interface{} `yaml:"payloads,omitempty" jsonschema:"title=payloads for the http request,description=Payloads contains any payloads for the current request"`

	// description: |
	//   Headers contains HTTP Headers to send with the request.
	// examples:
	//   - value: |
	//       map[string]string{"Content-Type": "application/x-www-form-urlencoded", "Content-Length": "1", "Any-Header": "Any-Value"}
	Headers map[string]string `yaml:"headers,omitempty" jsonschema:"title=headers to send with the http request,description=Headers contains HTTP Headers to send with the request"`
	// description: |
	//   RaceCount is the number of times to send a request in Race Condition Attack.
	// examples:
	//   - name: Send a request 5 times
	//     value: "5"
	RaceNumberRequests int `yaml:"race_count,omitempty" jsonschema:"title=number of times to repeat request in race condition,description=Number of times to send a request in Race Condition Attack"`
	// description: |
	//   MaxRedirects is the maximum number of redirects that should be followed.
	// examples:
	//   - name: Follow up to 5 redirects
	//     value: "5"
	MaxRedirects int `yaml:"max-redirects,omitempty" jsonschema:"title=maximum number of redirects to follow,description=Maximum number of redirects that should be followed"`
	// description: |
	//   PipelineConcurrentConnections is number of connections to create during pipelining.
	// examples:
	//   - name: Create 40 concurrent connections
	//     value: 40
	PipelineConcurrentConnections int `yaml:"pipeline-concurrent-connections,omitempty" jsonschema:"title=number of pipelining connections,description=Number of connections to create during pipelining"`
	// description: |
	//   PipelineRequestsPerConnection is number of requests to send per connection when pipelining.
	// examples:
	//   - name: Send 100 requests per pipeline connection
	//     value: 100
	PipelineRequestsPerConnection int `yaml:"pipeline-requests-per-connection,omitempty" jsonschema:"title=number of requests to send per pipelining connections,description=Number of requests to send per connection when pipelining"`
	// description: |
	//   Threads specifies number of threads to use sending requests. This enables Connection Pooling.
	//
	//   Connection: Close attribute must not be used in request while using threads flag, otherwise
	//   pooling will fail and engine will continue to close connections after requests.
	// examples:
	//   - name: Send requests using 10 concurrent threads
	//     value: 10
	Threads int `yaml:"threads,omitempty" jsonschema:"title=threads for sending requests,description=Threads specifies number of threads to use sending requests. This enables Connection Pooling"`
	// description: |
	//   MaxSize is the maximum size of http response body to read in bytes.
	// examples:
	//   - name: Read max 2048 bytes of the response
	//     value: 2048
	MaxSize       int  `yaml:"max-size,omitempty" jsonschema:"title=maximum http response body size,description=Maximum size of http response body to read in bytes"`
	SelfContained bool `yaml:"-" json:"-"`
	// description: |
	//   Signature is the request signature method
	// values:
	//   - "AWS"
	// description: |
	//   CookieReuse is an optional setting that enables cookie reuse for
	//   all requests defined in raw section.
	CookieReuse bool `yaml:"cookie-reuse,omitempty" jsonschema:"title=optional cookie reuse enable,description=Optional setting that enables cookie reuse"`
	// description: |
	//   Enables force reading of the entire raw unsafe request body ignoring
	//   any specified content length headers.
	ForceReadAllBody bool `yaml:"read-all,omitempty" jsonschema:"title=force read all body,description=Enables force reading of entire unsafe http request body"`
	// description: |
	//   Redirects specifies whether redirects should be followed by the HTTP Client.
	//
	//   This can be used in conjunction with `max-redirects` to control the HTTP request redirects.
	Redirects bool `yaml:"redirects,omitempty" jsonschema:"title=follow http redirects,description=Specifies whether redirects should be followed by the HTTP Client"`
	// description: |
	//   Pipeline defines if the attack should be performed with HTTP 1.1 Pipelining
	//
	//   All requests must be idempotent (GET/POST). This can be used for race conditions/billions requests.
	Pipeline bool `yaml:"pipeline,omitempty" jsonschema:"title=perform HTTP 1.1 pipelining,description=Pipeline defines if the attack should be performed with HTTP 1.1 Pipelining"`
	// description: |
	//   Unsafe specifies whether to use rawhttp engine for sending Non RFC-Compliant requests.
	//
	//   This uses the [rawhttp](https://github.com/projectdiscovery/rawhttp) engine to achieve complete
	//   control over the request, with no normalization performed by the client.
	Unsafe bool `yaml:"unsafe,omitempty" jsonschema:"title=use rawhttp non-strict-rfc client,description=Unsafe specifies whether to use rawhttp engine for sending Non RFC-Compliant requests"`
	// description: |
	//   Race determines if all the request have to be attempted at the same time (Race Condition)
	//
	//   The actual number of requests that will be sent is determined by the `race_count`  field.
	Race bool `yaml:"race,omitempty" jsonschema:"title=perform race-http request coordination attack,description=Race determines if all the request have to be attempted at the same time (Race Condition)"`
	// description: |
	//   ReqCondition automatically assigns numbers to requests and preserves their history.
	//
	//   This allows matching on them later for multi-request conditions.
	ReqCondition bool `yaml:"req-condition,omitempty" jsonschema:"title=preserve request history,description=Automatically assigns numbers to requests and preserves their history"`
	// description: |
	//   StopAtFirstMatch stops the execution of the requests and template as soon as a match is found.
	StopAtFirstMatch bool `yaml:"stop-at-first-match,omitempty" jsonschema:"title=stop at first match,description=Stop the execution after a match is found"`
	// description: |
	//   SkipVariablesCheck skips the check for unresolved variables in request
	SkipVariablesCheck bool `yaml:"skip-variables-check,omitempty" jsonschema:"title=skip variable checks,description=Skips the check for unresolved variables in request"`
	// description: |
	//   IterateAll iterates all the values extracted from internal extractors
	IterateAll bool `yaml:"iterate-all,omitempty" jsonschema:"title=iterate all the values,description=Iterates all the values extracted from internal extractors"`
	// description: |
	//   DigestAuthUsername specifies the username for digest authentication
	DigestAuthUsername string `yaml:"digest-username,omitempty" jsonschema:"title=specifies the username for digest authentication,description=Optional parameter which specifies the username for digest auth"`
	// description: |
	//   DigestAuthPassword specifies the password for digest authentication
	DigestAuthPassword string `yaml:"digest-password,omitempty" jsonschema:"title=specifies the password for digest authentication,description=Optional parameter which specifies the password for digest auth"`
}
