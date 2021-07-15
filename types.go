package main

import (
	"github.com/gastrodon/groudon/v2"
)

type CreateReportBody struct {
	Reported string `json:"reported"`
	Type     string `json:"type"`
	Reason   string `json:"reason"`
}

func (_ CreateReportBody) Validators() (values map[string]func(interface{}) (bool, error)) {
	values = map[string]func(interface{}) (bool, error){
		"reported": groudon.ValidUUID,
		"type":     groudon.ValidString,
		"reason":   groudon.ValidString,
	}

	return
}

func (_ CreateReportBody) Defaults() (values map[string]interface{}) {
	values = map[string]interface{}{
		"reason": "",
	}

	return
}

type PatchReportBody struct {
	Resolved   *bool   `json:"resolved"`
	Resolution *string `json:"resolution"`
}

func (_ PatchReportBody) Validators() (values map[string]func(interface{}) (bool, error)) {
	values = map[string]func(interface{}) (bool, error){
		"resolved":   groudon.OptionalBool,
		"resolution": groudon.OptionalString,
	}

	return
}

func (_ PatchReportBody) Defaults() (values map[string]interface{}) {
	values = map[string]interface{}{
		"resolution": "",
		"resolved":   nil,
	}

	return
}
