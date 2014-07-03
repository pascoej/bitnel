# bitnel/api/validator

At least it works. Written so that validation may be flexible.

Usage:

    type pers struct {
        Name  string
        Email string
    }

    func (a *pers) Rules() map[string][]Rule {
        return map[string][]Rule{
            "Name":  []Rule{&NonZero{}, &Length{3, 25}},
            "Email": []Rule{&NonZero{}, &Email{}},
        }
    }

    func TestValidation(t *testing.T) {
        asdf := &pers{
            Name:  "John Doe",
            Email: "andrewtidan.io",
        }

        v := New(asdf)

        if ok, errs := v.Validate(nil); !ok {
            fmt.Printf("%d errors received", len(errs))
            if !(len(errs) > 0) {
                t.Error("expecting errors")
            }
        }
    }