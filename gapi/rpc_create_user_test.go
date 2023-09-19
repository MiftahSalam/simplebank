package gapi

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	mockdb "simplebank/db/mock"
	db "simplebank/db/sqlc"
	"simplebank/pb"
	"simplebank/util"
	"simplebank/worker"
	mockwk "simplebank/worker/mock"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TestCaseUser struct {
	name          string
	body          *pb.CreateUserRequest
	buildStubs    func(store *mockdb.MockStore, worker *mockwk.MockTaskDistributor)
	checkResponse func(t *testing.T, res *pb.CreateUserResponse, err error)
}

type eqCreateUserParamsTxMatcher struct {
	arg      db.CreateUserTxParam
	password string
	user     db.User
}

func (expected eqCreateUserParamsTxMatcher) Matches(x interface{}) bool {
	actual, ok := x.(db.CreateUserTxParam)
	if !ok {
		return false
	}

	err := util.CheckPassword(expected.password, actual.HashedPassword)
	if err != nil {
		return false
	}

	expected.arg.HashedPassword = actual.HashedPassword

	if !reflect.DeepEqual(expected.arg.CreateUserParams, actual.CreateUserParams) {
		return false
	}

	err = actual.AfterCreate(expected.user)

	return err == nil
}

func (e eqCreateUserParamsTxMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParamsTx(arg db.CreateUserTxParam, password string, user db.User) gomock.Matcher {
	return eqCreateUserParamsTxMatcher{arg: arg, password: password, user: user}
}

func TestCreateUser(t *testing.T) {
	user, password := createRandomUser(t)

	testCases := []TestCaseUser{
		{
			name: "OK",
			body: &pb.CreateUserRequest{
				Username: user.Username,
				FullName: user.FullName,
				Email:    user.Email,
				Password: password,
			},
			buildStubs: func(store *mockdb.MockStore, task *mockwk.MockTaskDistributor) {
				arg := db.CreateUserTxParam{
					CreateUserParams: db.CreateUserParams{
						Username: user.Username,
						FullName: user.FullName,
						Email:    user.Email,
					},
				}
				store.EXPECT().CreateUserTx(gomock.Any(), EqCreateUserParamsTx(arg, password, user)).Times(1).Return(db.CreateUserTxResult{
					User: user,
				}, nil)

				payload := &worker.PayloadSendVerifyEmail{
					Username: user.Username,
				}
				task.EXPECT().DistributeTaskSendVerifyEmail(gomock.Any(), payload, gomock.Any()).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, user.Email, res.GetUser().Email)
				require.Equal(t, user.FullName, res.GetUser().FullName)
				require.Equal(t, user.Username, res.GetUser().Username)
			},
		},
		{
			name: "Internal",
			body: &pb.CreateUserRequest{
				Username: user.Username,
				FullName: user.FullName,
				Email:    user.Email,
				Password: password,
			},
			buildStubs: func(store *mockdb.MockStore, task *mockwk.MockTaskDistributor) {
				store.EXPECT().CreateUserTx(gomock.Any(), gomock.Any()).Times(1).Return(db.CreateUserTxResult{}, sql.ErrConnDone)

				task.EXPECT().DistributeTaskSendVerifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.Error(t, err)

				s, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, s.Code())
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storeCtrl := gomock.NewController(t)
			defer storeCtrl.Finish()

			workerCtrl := gomock.NewController(t)
			defer workerCtrl.Finish()

			store := mockdb.NewMockStore(storeCtrl)
			worker := mockwk.NewMockTaskDistributor(workerCtrl)

			tc.buildStubs(store, worker)

			server := newTestServer(t, store, worker)

			cur, err := server.CreateUser(context.Background(), tc.body)

			tc.checkResponse(t, cur, err)
		})
	}

}

func createRandomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	return
}
