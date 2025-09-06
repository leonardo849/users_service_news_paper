package units_test

import (
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
	"unicode"
	"users-service/config"
	"users-service/internal/helper"
	"users-service/pkg/date"
	"users-service/pkg/hash"
	"users-service/pkg/random"

	"github.com/google/uuid"
)

func TestMain(m *testing.M) {
	err := config.SetupEnvVar()
	if err != nil {
		log.Panic(err.Error())
	}
	code := m.Run()
	os.Exit(code)
}

func TestJwt(t *testing.T) {
	jwt, err := helper.GenerateJWT(uuid.NewString(),  time.Now(), helper.Customer)
	if err != nil {
		t.Error(err.Error())
	}
	parts := len(strings.Split(jwt, "."))
	if parts != 3  {
		t.Error("it isn't a jwt")
	}
}

func TestRoles(t *testing.T) {
	if !(helper.Ceo == "CEO" && helper.Customer == "CUSTOMER" && helper.Journalist == "JOURNALIST" && helper.Developer == "DEVELOPER" && len(helper.Roles) == 4) {
		t.Error("error in roles")
	}
} 

func TestHash(t *testing.T) {
	password := "batman123"
	hashP, err := hash.StringToHash(password)
	if err != nil {
		t.Error(err.Error())
	}
	if !hash.CompareHash(password, hashP) {
		t.Error("compare hash error")
	} 
	if hash.CompareHash(password+"32", hashP) {
		t.Error("compare hash error")
	}
}


func isPointer(i any) bool {
    return reflect.TypeOf(i).Kind() == reflect.Ptr
}




func TestGetPtrDate(t *testing.T) {
	now := time.Now()
	ptr := date.PtrTime(now)
	if !isPointer(ptr) {
		t.Error("it is not a pointer")
	}
}

func isNumeric(s string) bool {
    if s == "" {
        return false 
    }
    for _, r := range s {
        if !unicode.IsDigit(r) {
            return false
        }
    }
    return true
}

func TestRandomCode(t  *testing.T) {
	for i := 0; i<4; i++ {
		code := random.EncodeToString(6)
		if len(code) != 6 || !isNumeric(code) {
			t.Error("it doesn't have 6 characterics or it isn't numeric")
		}
	}
}