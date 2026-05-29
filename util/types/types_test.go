package types

import "testing"

func TestError_Struct(t *testing.T) {
	e := Error{Status: 400, Message: "Bad Request"}
	if e.Status != 400 {
		t.Errorf("expected status 400, got %d", e.Status)
	}
	if e.Message != "Bad Request" {
		t.Errorf("expected message 'Bad Request', got '%s'", e.Message)
	}
}

func TestOthers_AppendError(t *testing.T) {
	o := Others{}
	o.AppendError(Error{Status: 400, Message: "Error 1"})
	o.AppendError(Error{Status: 404, Message: "Error 2"})
	if len(o.Others) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(o.Others))
	}
	if o.Others[0].Message != "Error 1" {
		t.Errorf("expected 'Error 1', got '%s'", o.Others[0].Message)
	}
	if o.Others[1].Message != "Error 2" {
		t.Errorf("expected 'Error 2', got '%s'", o.Others[1].Message)
	}
}

func TestOthers_GetErrors_InvalidLimit(t *testing.T) {
	o := Others{}
	errs := o.GetErrors("", "", "0", "abc")
	if len(errs) == 0 {
		t.Error("expected errors for invalid limit")
	}
}

func TestOthers_GetErrors_NegativeLimit(t *testing.T) {
	o := Others{}
	errs := o.GetErrors("", "", "0", "-1")
	hasNegLimitErr := false
	for _, e := range errs {
		if e.Message == "gameId must be greater than 0" {
			hasNegLimitErr = true
		}
	}
	if !hasNegLimitErr {
		t.Error("expected error for negative limit")
	}
}

func TestOthers_GetErrors_InvalidPage(t *testing.T) {
	o := Others{}
	errs := o.GetErrors("", "", "abc", "10")
	hasInvalidPageErr := false
	for _, e := range errs {
		if e.Message == "Invalid page number" {
			hasInvalidPageErr = true
		}
	}
	if !hasInvalidPageErr {
		t.Error("expected error for invalid page")
	}
}

func TestOthers_GetErrors_NegativePage(t *testing.T) {
	o := Others{}
	errs := o.GetErrors("", "", "-1", "10")
	hasNegPageErr := false
	for _, e := range errs {
		if e.Message == "page less than 0" {
			hasNegPageErr = true
		}
	}
	if !hasNegPageErr {
		t.Error("expected error for negative page")
	}
}

func TestOthers_GetErrors_EmptyPage(t *testing.T) {
	o := Others{}
	errs := o.GetErrors("", "", "", "10")
	hasEmptyPageErr := false
	for _, e := range errs {
		if e.Message == "Page is required" {
			hasEmptyPageErr = true
		}
	}
	if !hasEmptyPageErr {
		t.Error("expected error for empty page")
	}
}

func TestOthers_GetErrors_EmptyLimit(t *testing.T) {
	o := Others{}
	errs := o.GetErrors("", "", "0", "")
	hasEmptyLimitErr := false
	for _, e := range errs {
		if e.Message == "Limit is required" {
			hasEmptyLimitErr = true
		}
	}
	if !hasEmptyLimitErr {
		t.Error("expected error for empty limit")
	}
}

func TestResponse_Struct(t *testing.T) {
	r := Response{Status: 200, Message: "OK", Data: "some data", Others: []Error{}}
	if r.Status != 200 {
		t.Errorf("expected status 200, got %d", r.Status)
	}
	if r.Message != "OK" {
		t.Errorf("expected message 'OK', got '%s'", r.Message)
	}
}

func TestData_Struct(t *testing.T) {
	p := Players{Players: []string{"P1", "P2"}}
	if len(p.Players) != 2 {
		t.Errorf("expected 2 players, got %d", len(p.Players))
	}
}

func TestEnemies_Struct(t *testing.T) {
	e := Enemies{Enemies: []string{"E1"}}
	if len(e.Enemies) != 1 {
		t.Errorf("expected 1 enemy, got %d", len(e.Enemies))
	}
}
