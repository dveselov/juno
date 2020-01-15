package providers

import "testing"

func TestUCallerSendSuccess(t *testing.T) {
	provider := UCallerProvider{
		ServiceID:  "517117",
		SecretKey:  "DzaNpa5yheif94H5stc8geRs7SRdkg6N",
		BaseAPIUrl: "https://api.ucaller.ru/v1.0",
	}
	err := provider.Send("79818246403", "1337")
	if err != nil {
		t.Fatal(err)
	}
}
