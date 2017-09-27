package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"
import "fmt"

const UnknownSaveErr = Error("UnknownSave")
const TranslationErr = Error("Translation")
const InvalidTUErr = Error("InvalidTU")

func convertSaveErrorCode(ec C.enum_CXSaveError) error {
	switch ec {
	case C.CXSaveError_None:
		return nil
	case C.CXSaveError_Unknown:
		return UnknownSaveErr
	case C.CXSaveError_TranslationErrors:
		return TranslationErr
	case C.CXSaveError_InvalidTU:
		return InvalidTUErr
	}

	return fmt.Errorf("unknown CXSaveError %d", int(ec))
}
