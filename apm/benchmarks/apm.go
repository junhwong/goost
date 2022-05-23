package benchmarks

import (
	"errors"
	"fmt"
	"time"

	"github.com/junhwong/goost/pkg/field"
)

var (
	errExample = errors.New("fail")

	_messages   = fakeMessages(1000)
	_tenInts    = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	_tenStrings = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	_tenTimes   = []time.Time{
		time.Unix(0, 0),
		time.Unix(1, 0),
		time.Unix(2, 0),
		time.Unix(3, 0),
		time.Unix(4, 0),
		time.Unix(5, 0),
		time.Unix(6, 0),
		time.Unix(7, 0),
		time.Unix(8, 0),
		time.Unix(9, 0),
	}
	_oneUser = &user{
		Name:      "Jane Doe",
		Email:     "jane@test.com",
		CreatedAt: time.Date(1980, 1, 1, 12, 0, 0, 0, time.UTC),
	}
	_tenUsers = users{
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
		_oneUser,
	}
)

func fakeMessages(n int) []string {
	messages := make([]string, n)
	for i := range messages {
		messages[i] = fmt.Sprintf("Test logging, but use a somewhat realistic message length. (#%v)", i)
	}
	return messages
}

func getMessage(iter int) string {
	return _messages[iter%1000]
}

type users []*user

// func (uu users) MarshalLogArray(arr zapcore.ArrayEncoder) error {
// 	var err error
// 	for i := range uu {
// 		err = multierr.Append(err, arr.AppendObject(uu[i]))
// 	}
// 	return err
// }

type user struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// func (u *user) MarshalLogObject(enc zapcore.ObjectEncoder) error {
// 	enc.AddString("name", u.Name)
// 	enc.AddString("email", u.Email)
// 	enc.AddInt64("createdAt", u.CreatedAt.UnixNano())
// 	return nil
// }

// func newZapLogger(lvl zapcore.Level) *zap.Logger {
// 	ec := zap.NewProductionEncoderConfig()
// 	ec.EncodeDuration = zapcore.NanosDurationEncoder
// 	ec.EncodeTime = zapcore.EpochNanosTimeEncoder
// 	enc := zapcore.NewJSONEncoder(ec)
// 	return zap.New(zapcore.NewCore(
// 		enc,
// 		&ztest.Discarder{},
// 		lvl,
// 	))
// }

// func newSampledLogger(lvl zapcore.Level) *zap.Logger {
// 	return zap.New(zapcore.NewSamplerWithOptions(
// 		newZapLogger(zap.DebugLevel).Core(),
// 		100*time.Millisecond,
// 		10, // first
// 		10, // thereafter
// 	))
// }

func fakeFields() []field.Field {
	_, i1 := field.Int("int")
	_, i2 := field.Slice("ints", field.IntKind)
	return []field.Field{
		i1(_tenInts[0]),
		i2(_tenInts),
		// field.String("string", _tenStrings[0]),
		// field.Strings("strings", _tenStrings),
		// field.Time("time", _tenTimes[0]),
		// field.Times("times", _tenTimes),
		// field.Object("user1", _oneUser),
		// field.Object("user2", _oneUser),
		// field.Array("users", _tenUsers),
		// field.Error(errExample),
	}
}

func fakeSugarFields() []interface{} {
	return []interface{}{
		"int", _tenInts[0],
		"ints", _tenInts,
		"string", _tenStrings[0],
		"strings", _tenStrings,
		"time", _tenTimes[0],
		"times", _tenTimes,
		"user1", _oneUser,
		"user2", _oneUser,
		"users", _tenUsers,
		"error", errExample,
	}
}

func fakeFmtArgs() []interface{} {
	// Need to keep this a function instead of a package-global var so that we
	// pay the cast-to-interface{} penalty on each call.
	return []interface{}{
		_tenInts[0],
		_tenInts,
		_tenStrings[0],
		_tenStrings,
		_tenTimes[0],
		_tenTimes,
		_oneUser,
		_oneUser,
		_tenUsers,
		errExample,
	}
}
