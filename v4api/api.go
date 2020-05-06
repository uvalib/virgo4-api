package v4api

// SearchResponse contains teh aggregated and sorted search response data from all pools in the system
type SearchResponse struct {
	Request     *SearchRequest `json:"request"`
	Pools       []PoolIdentity `json:"pools"`
	TotalTimeMS int64          `json:"total_time_ms"`
	TotalHits   int            `json:"total_hits"`
	Results     []*PoolResult  `json:"pool_results"`
	Warnings    []string       `json:"warnings"`
	Suggestions []Suggestion   `json:"suggestions"`
}

// Suggestion contains search suggestion data
type Suggestion struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// SortOptionEnum is the enumerated type for WorldCat sort options
type SortOptionEnum int

const (
	// SortRelevance is used to sort by descending relevance
	SortRelevance SortOptionEnum = iota
	// SortDate is used to sort by published date
	SortDate
	// SortTitle is used to sort by title
	SortTitle
	// SortAuthor is used to sort by Author
	SortAuthor
)

func (r SortOptionEnum) String() string {
	return []string{"SortRelevance", "SortDatePublished", "SortTitle", "SortAuthor"}[r]
}

// PoolAttribute describes a capability of a pool
type PoolAttribute struct {
	Name      string `json:"name"`
	Supported bool   `json:"supported"`
	Value     string `json:"value,omitempty"`
}

// PoolIdentity contains the complete data that descibes a pool and its abilities
type PoolIdentity struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Mode        string          `json:"mode"`
	URL         string          `json:"url"`
	Attributes  []PoolAttribute `json:"attributes,omitempty"`
	SortOptions []SortOption    `json:"sort_options,omitempty"`
}

// SortOrder specifies sort ordering for a given search.
type SortOrder struct {
	SortID string `json:"sort_id"`
	Order  string `json:"order"`
}

// SortOption defines a sorting option for a pool
type SortOption struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// SearchRequest contains all of the data necessary for a client seatch request
type SearchRequest struct {
	Query       string            `json:"query"`
	Pagination  Pagination        `json:"pagination"`
	Sort        SortOrder         `json:"sort,omitempty"`
	Filters     []Filter          `json:"filters,omitempty"`
	Preferences SearchPreferences `json:"preferences,omitempty"`
}

// Pagination cantains pagination info
type Pagination struct {
	Start int `json:"start"`
	Rows  int `json:"rows"`
	Total int `json:"total"`
}

// Filter contains the fields for a single filter.
type Filter struct {
	PoolID string `json:"pool_id"`
	Facets []struct {
		FacetID string `json:"facet_id"`
		Value   string `json:"value"`
	} `json:"facets"`
}

// PoolResult contains search responce data from a pool
type PoolResult struct {
	ServiceURL      string                 `json:"service_url,omitempty"`
	PoolName        string                 `json:"pool_id,omitempty"`
	ElapsedMS       int64                  `json:"elapsed_ms,omitempty"`
	Pagination      Pagination             `json:"pagination"`
	Sort            SortOrder              `json:"sort,omitempty"`
	Groups          []Group                `json:"group_list,omitempty"`
	FacetList       []Facet                `json:"facet_list,omitempty"`
	Confidence      string                 `json:"confidence,omitempty"`
	Debug           map[string]interface{} `json:"debug"`
	Warnings        []string               `json:"warnings"`
	StatusCode      int                    `json:"status_code"`
	StatusMessage   string                 `json:"status_msg,omitempty"`
	ContentLanguage string                 `json:"-"`
}

// ConfidenceIndex will convert a string confidence into a numeric value
// with low having the lowest value and exact the highest
func (pr *PoolResult) ConfidenceIndex() int {
	conf := []string{"low", "medium", "high", "exact"}
	for idx, val := range conf {
		if val == pr.Confidence {
			return idx
		}
	}
	// No confidence match. Assume worst value
	return 0
}

// Facet contains the fields for a single facet.
type Facet struct {
	ID      string        `json:"id"`
	Name    string        `json:"name"`
	Type    string        `json:"type"`
	Buckets []FacetBucket `json:"buckets,omitempty"`
}

// FacetBucket contains the fields for an individual bucket for a facet.
type FacetBucket struct {
	Value    string `json:"value"`
	Count    int    `json:"count"`
	Selected bool   `json:"selected"`
}

// Group contains the records for a single group in a search result set.
type Group struct {
	Value   string   `json:"value"`
	Count   int      `json:"count"`
	Records []Record `json:"record_list,omitempty"`
}

// RelatedRecord contains basic info for other records with the same
// group value (currently only used by the Solr image pool).
type RelatedRecord struct {
	ID              string `json:"id,omitempty"`
	IIIFManifestURL string `json:"iiif_manifest_url,omitempty"`
	IIIFImageURL    string `json:"iiif_image_url,omitempty"`
	IIIFBaseURL     string `json:"iiif_base_url,omitempty"`
}

// Record is a summary of one search hit
type Record struct {
	Fields     []RecordField          `json:"fields"`
	Related    []RelatedRecord        `json:"related,omitempty"`
	Debug      map[string]interface{} `json:"debug"`
	GroupValue string                 `json:"-"` // used in Solr pools to properly group results
}

// RecordField contains metadata for a single field in a record.
type RecordField struct {
	Name       string `json:"name"`
	Type       string `json:"type,omitempty"` // empty implies "text"
	Label      string `json:"label"`
	Value      string `json:"value"`
	Visibility string `json:"visibility,omitempty"` // e.g. "basic" or "detailed".  empty implies "basic"
	Display    string `json:"display,omitempty"`    // e.g. "optional".  empty implies not optional
	Provider   string `json:"provider,omitempty"`   // for URLs (e.g. "hathitrust", "proquest")
	Item       string `json:"item,omitempty"`       // for certain URLs (currently hathitrust)
	RISCode    string `json:"ris_code,omitempty"`
}

// SearchPreferences contains preferences for the search
type SearchPreferences struct {
	TargetPool   string   `json:"target_pool"`
	ExcludePools []string `json:"exclude_pool"`
}

// IsExcluded will return true if the target URL is included in the ExcludePools preferece
func (p *SearchPreferences) IsExcluded(URL string) bool {
	if URL == "" {
		return false
	}
	for _, excludedURL := range p.ExcludePools {
		if excludedURL == URL {
			return true
		}
	}
	return false
}
