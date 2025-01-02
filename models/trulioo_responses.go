package models

type InitTrulioo struct {
	CanGoBack bool `json:"canGoBack"`
	Elements  []struct {
		ID             string `json:"id"`
		Name           string `json:"name"`
		PrefillOptions string `json:"prefillOptions,omitempty"`
		Role           string `json:"role,omitempty"`
		Type           string `json:"type"`
		Validations    []struct {
			Type string `json:"type"`
		} `json:"validations,omitempty"`
		Condition struct {
			Rules []interface{} `json:"rules"`
			Type  string        `json:"type"`
		} `json:"condition,omitempty"`
		NormalizedName string `json:"normalizedName,omitempty"`
		DateFormat     string `json:"dateFormat,omitempty"`
		Content        string `json:"content,omitempty"`
		Placeholder    string `json:"placeholder,omitempty"`
		Value          string `json:"value,omitempty"`
	} `json:"elements"`
	ID       string `json:"id"`
	Subtitle string `json:"subtitle"`
	Title    string `json:"title"`
	Type     string `json:"type"`
}

type DirectSubmitResponse struct {
	ID          string `json:"id"`
	RedirectURL string `json:"redirectUrl"`
	Text        string `json:"text"`
	Type        string `json:"type"`
}

type BearerTokenReponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type ClientDetailsResponse struct {
	ID                    string        `json:"id"`
	OwnerIds              []interface{} `json:"ownerIds"`
	Created               int           `json:"created"`
	LastModified          int           `json:"lastModified"`
	ProfileType           string        `json:"profileType"`
	Status                string        `json:"status"`
	StatusManuallyChanged bool          `json:"statusManuallyChanged"`
	FlowData              struct {
		Six74E0Ec2678Bea7Ec730Fd98 struct {
			ID        string `json:"id"`
			Completed bool   `json:"completed"`
			FieldData struct {
				Six716B75A1287D277472C8D84 struct {
					ID    string   `json:"id"`
					Name  string   `json:"name"`
					Value []string `json:"value"`
					Role  string   `json:"role"`
				} `json:"6716b75a1287d277472c8d84"`
				Six7228Aef1E5E2108D84020A2 struct {
					ID    string   `json:"id"`
					Name  string   `json:"name"`
					Value []string `json:"value"`
					Role  string   `json:"role"`
				} `json:"67228aef1e5e2108d84020a2"`
				Six716B75A1287D277472C8D81 struct {
					ID    string   `json:"id"`
					Name  string   `json:"name"`
					Value []string `json:"value"`
					Role  string   `json:"role"`
				} `json:"6716b75a1287d277472c8d81"`
				Six716B75A1287D277472C8D87 struct {
					ID    string   `json:"id"`
					Name  string   `json:"name"`
					Value []string `json:"value"`
					Role  string   `json:"role"`
				} `json:"6716b75a1287d277472c8d87"`
				Six716B75A1287D277472C8D8C struct {
					ID    string   `json:"id"`
					Name  string   `json:"name"`
					Value []string `json:"value"`
					Role  string   `json:"role"`
				} `json:"6716b75a1287d277472c8d8c"`
				Six716B75A1287D277472C8D82 struct {
					ID    string   `json:"id"`
					Name  string   `json:"name"`
					Value []string `json:"value"`
					Role  string   `json:"role"`
				} `json:"6716b75a1287d277472c8d82"`
				Six716B75A1287D277472C8D86 struct {
					ID    string   `json:"id"`
					Name  string   `json:"name"`
					Value []string `json:"value"`
					Role  string   `json:"role"`
				} `json:"6716b75a1287d277472c8d86"`
				Six744Facf99661447B4B58Ff7 struct {
					ID    string   `json:"id"`
					Name  string   `json:"name"`
					Value []string `json:"value"`
					Role  string   `json:"role"`
				} `json:"6744facf99661447b4b58ff7"`
				Six716B75A1287D277472C8D83 struct {
					ID    string   `json:"id"`
					Name  string   `json:"name"`
					Value []string `json:"value"`
					Role  string   `json:"role"`
				} `json:"6716b75a1287d277472c8d83"`
				Six716B75A1287D277472C8D88 struct {
					ID    string   `json:"id"`
					Name  string   `json:"name"`
					Value []string `json:"value"`
					Role  string   `json:"role"`
				} `json:"6716b75a1287d277472c8d88"`
				Six716B75A1287D277472C8D85 struct {
					ID             string   `json:"id"`
					Name           string   `json:"name"`
					Value          []string `json:"value"`
					NormalizedName string   `json:"normalizedName"`
				} `json:"6716b75a1287d277472c8d85"`
			} `json:"fieldData"`
			FileData struct {
			} `json:"fileData"`
			EnhancedFileData struct {
			} `json:"enhancedFileData"`
			ServiceData []struct {
				Timestamp       int    `json:"timestamp"`
				ServiceStatus   string `json:"serviceStatus"`
				NodeID          string `json:"nodeId"`
				NodeTitle       string `json:"nodeTitle"`
				NodeType        string `json:"nodeType"`
				Match           bool   `json:"match"`
				TransactionInfo struct {
					AccountName         string `json:"accountName"`
					AccountIdentifier   string `json:"accountIdentifier"`
					CountryCode         string `json:"countryCode"`
					Date                string `json:"date"`
					TransactionID       string `json:"transactionId"`
					TransactionRecordID string `json:"transactionRecordId"`
					UserName            string `json:"userName"`
				} `json:"transactionInfo"`
				FullServiceDetails struct {
					TransactionID string `json:"TransactionID"`
					UploadedDt    string `json:"UploadedDt"`
					CompletedDt   string `json:"CompletedDt"`
					CountryCode   string `json:"CountryCode"`
					ProductName   string `json:"ProductName"`
					Record        struct {
						TransactionRecordID string `json:"TransactionRecordID"`
						RecordStatus        string `json:"RecordStatus"`
						DatasourceResults   []struct {
							DatasourceName   string        `json:"DatasourceName"`
							DatasourceFields []interface{} `json:"DatasourceFields"`
							AppendedFields   []struct {
								FieldName string `json:"FieldName"`
								Data      string `json:"Data"`
							} `json:"AppendedFields"`
							Errors      []interface{} `json:"Errors"`
							FieldGroups []interface{} `json:"FieldGroups"`
						} `json:"DatasourceResults"`
						Errors []interface{} `json:"Errors"`
						Rule   struct {
							RuleName string `json:"RuleName"`
							Note     string `json:"Note"`
						} `json:"Rule"`
						SupplementaryRules []interface{} `json:"SupplementaryRules"`
					} `json:"Record"`
					Errors            []interface{} `json:"Errors"`
					UserGUID          string        `json:"UserGUID"`
					AccountIdentifier string        `json:"AccountIdentifier"`
					InputFields       []struct {
						FieldName string `json:"FieldName"`
						Value     string `json:"Value"`
					} `json:"InputFields"`
				} `json:"fullServiceDetails"`
				WatchlistResults struct {
					AdvancedWatchlist struct {
						WatchlistStatus     string `json:"watchlistStatus"`
						WatchlistHitDetails struct {
							WlHitsNumber  int `json:"wlHitsNumber"`
							AmHitsNumber  int `json:"amHitsNumber"`
							PepHitsNumber int `json:"pepHitsNumber"`
						} `json:"watchlistHitDetails"`
					} `json:"Advanced Watchlist"`
				} `json:"watchlistResults,omitempty"`
			} `json:"serviceData"`
		} `json:"674e0ec2678bea7ec730fd98"`
	} `json:"flowData"`
	Notes []interface{} `json:"notes"`
	Tasks []interface{} `json:"tasks"`
	User  string        `json:"user"`
}
