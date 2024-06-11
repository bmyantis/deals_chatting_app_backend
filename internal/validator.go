// validator.go
package validator

type Validator interface {
    Struct(interface{}) error
}
