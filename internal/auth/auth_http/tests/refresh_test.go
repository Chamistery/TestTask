package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/Chamistery/TestTask/internal/auth/auth_http"
	"github.com/Chamistery/TestTask/internal/auth/model"
	"github.com/Chamistery/TestTask/internal/auth/service/mocks"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefresh(t *testing.T) {
	ctx := context.Background()

	mc := minimock.NewController(t)
	defer mc.Finish()

	mockSvc := mocks.NewAuthServiceMock(mc)
	impl := auth_http.NewImplementation(mockSvc)

	mockSvc.CreateMock.Set(func(ctx context.Context, m model.CreateModel) (string, error) {
		return "created-guid", nil
	})

	createResp, err := impl.Create(ctx, &model.CreateRequest{Guid: "user-123"}, "127.0.0.1")
	require.NoError(t, err)
	validToken := createResp.RefreshToken
	origIP := "127.0.0.1"

	tests := []struct {
		name            string
		reqToken        string
		callIP          string
		mockRefreshFunc func(ctx context.Context, m model.RefreshModel) (string, error)
		expectError     bool
		expectedErrMsg  string
	}{
		{
			name:     "Success matching IP",
			reqToken: validToken,
			callIP:   origIP,
			mockRefreshFunc: func(ctx context.Context, m model.RefreshModel) (string, error) {
				if m.Uuid == "" || m.RefreshTokenOld != validToken {
					return "", errors.New("bad model")
				}
				return "returned-guid", nil
			},
			expectError:    false,
			expectedErrMsg: "",
		},
		{
			name:            "Invalid token string",
			reqToken:        "not-a-jwt",
			callIP:          origIP,
			mockRefreshFunc: nil,
			expectError:     true,
			expectedErrMsg:  "invalid refresh token",
		},
		{
			name:     "Service returns error",
			reqToken: validToken,
			callIP:   origIP,
			mockRefreshFunc: func(ctx context.Context, m model.RefreshModel) (string, error) {
				return "", errors.New("service error")
			},
			expectError:    true,
			expectedErrMsg: "service error",
		},
		{
			name:     "IP mismatch",
			reqToken: validToken,
			callIP:   "192.168.0.42",
			mockRefreshFunc: func(ctx context.Context, m model.RefreshModel) (string, error) {
				return "returned-guid", nil
			},
			expectError:    true,
			expectedErrMsg: "Ip changed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.mockRefreshFunc != nil {
				mockSvc.RefreshMock.Set(tc.mockRefreshFunc)
			}

			req := &model.RefreshRequest{RefreshToken: tc.reqToken}
			resp, err := impl.Refresh(ctx, req, tc.callIP)

			if tc.expectError {
				assert.Error(t, err)
				if tc.expectedErrMsg != "" {
					assert.Equal(t, tc.expectedErrMsg, err.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotEmpty(t, resp.AccessToken)
				assert.NotEmpty(t, resp.RefreshToken)
			}
		})
	}
}
