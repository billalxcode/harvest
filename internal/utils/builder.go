package utils

type URLBuilder struct {
	scheme string
	domain string
	port   string
	path   string
	params map[string]string
}

func NewURLBuilder(scheme string, domain string, port string, path string, params map[string]string) *URLBuilder {
	if port == "" {
		if scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}
	
	return &URLBuilder{
		scheme: scheme,
		domain: domain,
		port:   port,
		path:   path,
		params: params,
	}
}

func (b *URLBuilder) SetScheme(scheme string) {
	b.scheme = scheme
}

func (b *URLBuilder) SetDomain(domain string) {
	b.domain = domain
}
func (b *URLBuilder) SetPort(port string) {
	b.port = port
}
func (b *URLBuilder) SetPath(path string) {
	b.path = path
}
func (b *URLBuilder) SetParams(params map[string]string) {
	b.params = params
}
func (b *URLBuilder) Build() string {
	url := b.scheme + "://" + b.domain + ":" + b.port + "/" + b.path

	if b.params != nil {
		url += "?"
		for key, value := range b.params {
			url += key + "=" + value
		}
	}

	return url
}
