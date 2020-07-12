package model

// CaddyJSONConfigModel 是caddy api的json配置文件
type CaddyJSONConfigModel struct {
	Admin struct {
		Disabled      bool     `json:"disabled"`
		Listen        string   `json:"listen"`
		EnforceOrigin bool     `json:"enforce_origin"`
		Origins       []string `json:"origins"`
		Config        struct {
			Persist bool `json:"persist"`
		} `json:"config"`
	} `json:"admin"`
	Logging struct {
		Sink struct {
			Writer struct {
			} `json:"writer"`
		} `json:"sink"`
		Logs struct {
			LogDetail struct {
				Writer struct {
				} `json:"writer"`
				Encoder struct {
				} `json:"encoder"`
				Level    string `json:"level"`
				Sampling struct {
					Interval   int `json:"interval"`
					First      int `json:"first"`
					Thereafter int `json:"thereafter"`
				} `json:"sampling"`
				Include []string `json:"include"`
				Exclude []string `json:"exclude"`
			} `json:"log_detail"`
		} `json:"logs"`
	} `json:"logging"`
	Storage struct {
	} `json:"storage"`
	Apps struct {
		HTTP struct {
			HTTPPort    int `json:"http_port"`
			HTTPSPort   int `json:"https_port"`
			GracePeriod int `json:"grace_period"`
			Servers     struct {
				CaddyWebGui struct {
					Listen           []string `json:"listen"`
					ListenerWrappers []struct {
					} `json:"listener_wrappers"`
					ReadTimeout       int `json:"read_timeout"`
					ReadHeaderTimeout int `json:"read_header_timeout"`
					WriteTimeout      int `json:"write_timeout"`
					IdleTimeout       int `json:"idle_timeout"`
					MaxHeaderBytes    int `json:"max_header_bytes"`
					Routes            []struct {
						Group string `json:"group"`
						Match []struct {
							Host string `json:"host"`
						} `json:"match"`
						Handle []struct {
						} `json:"handle"`
						Terminal bool `json:"terminal"`
					} `json:"routes"`
					Errors struct {
						Routes []struct {
							Group string `json:"group"`
							Match []struct {
							} `json:"match"`
							Handle []struct {
							} `json:"handle"`
							Terminal bool `json:"terminal"`
						} `json:"routes"`
					} `json:"errors"`
					TLSConnectionPolicies []struct {
						Match struct {
						} `json:"match"`
						CertificateSelection struct {
							SerialNumber []struct {
							} `json:"serial_number"`
							SubjectOrganization []string `json:"subject_organization"`
							PublicKeyAlgorithm  int      `json:"public_key_algorithm"`
							AnyTag              []string `json:"any_tag"`
							AllTags             []string `json:"all_tags"`
						} `json:"certificate_selection"`
						CipherSuites         []string `json:"cipher_suites"`
						Curves               []string `json:"curves"`
						Alpn                 []string `json:"alpn"`
						ProtocolMin          string   `json:"protocol_min"`
						ProtocolMax          string   `json:"protocol_max"`
						ClientAuthentication struct {
							TrustedCaCerts   []string `json:"trusted_ca_certs"`
							TrustedLeafCerts []string `json:"trusted_leaf_certs"`
							Mode             string   `json:"mode"`
						} `json:"client_authentication"`
						DefaultSni string `json:"default_sni"`
					} `json:"tls_connection_policies"`
					AutomaticHTTPS struct {
						Disable                  bool     `json:"disable"`
						DisableRedirects         bool     `json:"disable_redirects"`
						Skip                     []string `json:"skip"`
						SkipCertificates         []string `json:"skip_certificates"`
						IgnoreLoadedCertificates bool     `json:"ignore_loaded_certificates"`
					} `json:"automatic_https"`
					StrictSniHost bool `json:"strict_sni_host"`
					Logs          struct {
						DefaultLoggerName string `json:"default_logger_name"`
						LoggerNames       struct {
							string `json:""`
						} `json:"logger_names"`
						SkipHosts         []string `json:"skip_hosts"`
						SkipUnmappedHosts bool     `json:"skip_unmapped_hosts"`
					} `json:"logs"`
					ExperimentalHTTP3 bool `json:"experimental_http3"`
					AllowH2C          bool `json:"allow_h2c"`
				} `json:"caddy_web_gui"`
			} `json:"servers"`
		} `json:"http"`
	} `json:"apps"`
}
