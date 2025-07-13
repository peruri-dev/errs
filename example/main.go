package main

import (
	"fmt"
	"net/http"

	"github.com/peruri-dev/errs"
)

func eOrigin() error {
	return errs.New("the original error")
}

func process1() error {
	e := eOrigin()
	if e != nil {
		return errs.Chain(e, "got problem while process1 call")
	}

	return nil
}

func process2() error {
	return errs.Chain(process1(), "processing 2 gain error also")
}

func usecase() error {
	err := process2()

	ErrInvalidParam := errs.NewCodex("Terjadi Kesalahan", "Harap lakukan langkah 1 dan 2", "000-123", http.StatusBadRequest)

	return errs.ChainCodex(err, ErrInvalidParam.SetErr("please hide this message for internal purpose only!"))
}

func main() {
	err := usecase()
	if err == nil {
		panic("should error")
	}

	eps := errs.PrintStack(err)
	epsj := errs.PrintStackJson(err)
	ec := errs.ParseCodex(err)

	fmt.Println("Print directly:", err)
	fmt.Println("Print error():", err.Error())
	fmt.Println("Print stack trace:", eps)
	fmt.Println("Print stack trace in json:", epsj)
	fmt.Printf("Get codex information: %+v\n", ec)
}
