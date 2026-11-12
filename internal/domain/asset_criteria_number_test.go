package domain

import "testing"

func TestValidateAssetCriteriaNumber(t *testing.T) {
	tests := []struct {
		raw   string
		valid bool
	}{
		{valid: true, raw: "9007199254740993"},
		{valid: true, raw: "1.25e2"},
		{valid: true, raw: "1e-6176"},
		{valid: true, raw: "10e-6177"},
		{valid: true, raw: "1e6144"},
		{valid: false, raw: "1e-6177"},
		{valid: false, raw: "1e6145"},
		{valid: false, raw: "12345678901234567890123456789012345"},
		{valid: false, raw: "1e7000"},
		{valid: false, raw: "null"},
		{valid: false, raw: `"1"`},
		{valid: false, raw: "+1"},
	}
	for _, test := range tests {
		t.Run(test.raw, func(t *testing.T) {
			err := ValidateAssetCriteriaNumber(test.raw)
			if (err == nil) != test.valid {
				t.Fatalf("ValidateAssetCriteriaNumber(%q) error = %v, valid = %t", test.raw, err, test.valid)
			}
		})
	}
}
