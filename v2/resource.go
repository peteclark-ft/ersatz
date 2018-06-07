package v2

import (
	"encoding/json"
	"net/textproto"
	"net/url"
)

func (r *Resource) UnmarshalJSON(d []byte) error {
	resp := &Response{}
	err := json.Unmarshal(d, resp)

	if err == nil {
		r.Response = *resp
		return nil
	}

	arr := make([]Discriminator, 0)
	err = json.Unmarshal(d, &arr)
	if err != nil {
		return err
	}

	r.Discriminators = append(r.Discriminators, arr...)
	return nil
}

// UnmarshalJSON creates a textproto.MIMEHeader compliant struct using provided map values
func (h *Headers) UnmarshalJSON(d []byte) error {
	headers := make(map[string]string)
	err := json.Unmarshal(d, &headers)
	if err != nil {
		return err
	}

	templated, remainder := ParseRequestValues(headers)
	h.TemplatedValues = templated

	h.MIMEHeader = textproto.MIMEHeader{}
	for k, v := range remainder {
		h.Add(k, v)
	}
	return nil
}

// UnmarshalJSON creates a url.Values compliant struct using the provided map values
func (q *QueryParams) UnmarshalJSON(d []byte) error {
	query := make(map[string]string)
	err := json.Unmarshal(d, &query)
	if err != nil {
		return err
	}

	templated, remainder := ParseRequestValues(query)
	q.TemplatedValues = templated

	q.Values = url.Values{}
	for k, v := range remainder {
		q.Add(k, v)
	}
	return nil
}
