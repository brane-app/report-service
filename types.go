package main

type CreateReportBody struct {
	Reported string `json:"reported"`
	Type     string `json:"type"`
	Reason   string `json:"reason"`
}

func (_ CreateReportBody) Types() (values map[string]string) {
	values = map[string]string{
		"reported": "string",
		"type":     "string",
		"reason":   "string",
	}

	return
}

func (_ CreateReportBody) Defaults() (values map[string]interface{}) {
	values = map[string]interface{}{
		"reason": "",
	}

	return
}
