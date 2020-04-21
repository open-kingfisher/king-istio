package chart

type Chart struct {
	Name      string                 `json:"name"`
	Value     map[string]string      `json:"value,omitempty"`
	Children  []Chart                `json:"children,omitempty"`
	LineStyle map[string]string      `json:"lineStyle,omitempty"`
	ItemStyle map[string]interface{} `json:"itemStyle,omitempty"`
	Rank      string                 `json:"rank,omitempty"`
	Count     int                    `json:"count,omitempty"`
}
