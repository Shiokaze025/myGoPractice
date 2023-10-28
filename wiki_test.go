package main

import "testing"

var testPage = Page{
	Title: "hello",
	Body:  []byte("hello world"),
}

func TestPage_save(t *testing.T) {
	type fields struct {
		Title string
		Body  []byte
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
		{"save", fields(testPage), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Page{
				Title: tt.fields.Title,
				Body:  tt.fields.Body,
			}
			if err := p.save(); (err != nil) != tt.wantErr {
				t.Errorf("save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
