package user

import (
	"bytes"
	"database/sql"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/benbjohnson/clock"
	"github.com/gofiber/fiber/v2"
	domain "github.com/solutionchallenge/ondaum-server/internal/domain/user"
	pkghttp "github.com/solutionchallenge/ondaum-server/pkg/http"
	"github.com/solutionchallenge/ondaum-server/pkg/jwt"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
	"github.com/solutionchallenge/ondaum-server/test/mock"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/schema"
	"go.uber.org/mock/gomock"
)

var testcases_GetSelfHandler = []struct {
	name     string
	setup    func(t *testing.T, tester *GetSelfHandlerTester)
	request  *http.Request
	response *http.Response
}{
	{
		name: "Failure Case - Unauthorized",
		setup: func(t *testing.T, tester *GetSelfHandlerTester) {
			tester.mockedJWT.EXPECT().GetTokenType(gomock.Any()).Return(jwt.InvalidType, nil)
		},
		request: &http.Request{
			Method: http.MethodGet,
			URL:    &url.URL{Path: "/"},
			Header: http.Header{
				"Authorization": []string{"Bearer invalid"},
			},
		},
		response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       io.NopCloser(strings.NewReader(`{"message":"Unauthorized"}`)),
		},
	},
	{
		name: "Failure Case - User Not Found",
		setup: func(t *testing.T, tester *GetSelfHandlerTester) {
			preparedClaims := jwt.Claims{
				Value: "1",
				Metadata: map[string]any{
					"test": "test",
				},
			}
			tester.mockedJWT.EXPECT().GetTokenType(gomock.Any()).Return(jwt.AccessTokenType, nil)
			tester.mockedJWT.EXPECT().UnpackToken(gomock.Any()).Return(&preparedClaims, nil)

			query := tester.mockedORM.NewSelect().
				Model((*domain.User)(nil)).
				Relation("Addition").
				Relation("Privacy").
				Where("id = ?", 1)
			queryString, _ := query.AppendQuery(schema.NewFormatter(tester.mockedORM.Dialect()), nil)
			escapedQueryString := regexp.QuoteMeta(string(queryString))
			tester.databaseController.
				ExpectQuery(escapedQueryString).
				WillReturnError(sql.ErrNoRows)
		},
		request: &http.Request{
			Method: http.MethodGet,
			URL:    &url.URL{Path: "/"},
			Header: http.Header{
				"Authorization": []string{"Bearer valid"},
			},
		},
		response: &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(strings.NewReader(`{"message":"User not found for id: 1"}`)),
		},
	},
	{
		name: "Failure Case - Internal Server Error",
		setup: func(t *testing.T, tester *GetSelfHandlerTester) {
			preparedClaims := jwt.Claims{
				Value: "1",
			}
			tester.mockedJWT.EXPECT().GetTokenType(gomock.Any()).Return(jwt.AccessTokenType, nil)
			tester.mockedJWT.EXPECT().UnpackToken(gomock.Any()).Return(&preparedClaims, nil)

			query := tester.mockedORM.NewSelect().
				Model((*domain.User)(nil)).
				Relation("Addition").
				Relation("Privacy").
				Where("id = ?", 1)
			queryString, _ := query.AppendQuery(schema.NewFormatter(tester.mockedORM.Dialect()), nil)
			escapedQueryString := regexp.QuoteMeta(string(queryString))
			tester.databaseController.
				ExpectQuery(escapedQueryString).
				WillReturnError(sql.ErrConnDone)
		},
		request: &http.Request{
			Method: http.MethodGet,
			URL:    &url.URL{Path: "/"},
			Header: http.Header{
				"Authorization": []string{"Bearer valid"},
			},
		},
		response: &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(strings.NewReader(`{"message":"Failed to get user for id: 1"}`)),
		},
	},
	{
		name: "Success Case",
		setup: func(t *testing.T, tester *GetSelfHandlerTester) {
			preparedClaims := jwt.Claims{
				Value: "1",
			}
			tester.mockedJWT.EXPECT().GetTokenType(gomock.Any()).Return(jwt.AccessTokenType, nil)
			tester.mockedJWT.EXPECT().UnpackToken(gomock.Any()).Return(&preparedClaims, nil)

			query := tester.mockedORM.NewSelect().
				Model((*domain.User)(nil)).
				Relation("Addition").
				Relation("Privacy").
				Where("id = ?", 1)
			queryString, _ := query.AppendQuery(schema.NewFormatter(tester.mockedORM.Dialect()), nil)
			escapedQueryString := regexp.QuoteMeta(string(queryString))
			tester.databaseController.
				ExpectQuery(escapedQueryString).
				WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "created_at", "updated_at"}).
					AddRow(
						1,
						"John Doe",
						"john.doe@example.com",
						time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
						time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					),
				)
		},
		request: &http.Request{
			Method: http.MethodGet,
			URL:    &url.URL{Path: "/"},
			Header: http.Header{
				"Authorization": []string{"Bearer valid"},
			},
		},
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(bytes.NewReader(utils.MustMarshal(domain.UserDTO{
				ID:       1,
				Username: "John Doe",
				Email:    "john.doe@example.com",
			}, utils.MarshalJSON))),
		},
	},
}

func Test_GetSelfHandler(t *testing.T) {
	tester, err := prepareGetSelfHandlerForTest(t)
	if err != nil {
		t.Fatalf("failed to prepare tester: %v", err)
	}

	app := fiber.New()
	app.Use(pkghttp.NewJWTAuthMiddleware(tester.mockedJWT))
	app.Get("/", tester.handler.Handle)

	for _, testcase := range testcases_GetSelfHandler {
		t.Run(testcase.name, func(t *testing.T) {
			testcase.setup(t, tester)
			response, err := app.Test(testcase.request)
			if err != nil {
				t.Fatalf("failed to test: %v", err)
			}
			if response.StatusCode != testcase.response.StatusCode {
				t.Fatalf("expected status code %d, got %d", testcase.response.StatusCode, response.StatusCode)
			}
			expectedBody, err := io.ReadAll(testcase.response.Body)
			if err != nil {
				t.Fatalf("failed to read expected response body: %v", err)
			}
			actualBody, err := io.ReadAll(response.Body)
			if err != nil {
				t.Fatalf("failed to read response body: %v", err)
			}
			if string(actualBody) != string(expectedBody) {
				t.Fatalf("expected body %s, got %s", expectedBody, actualBody)
			}
		})
	}

	tester.mockedORM.Close()
	tester.mockedDatabase.Close()
	tester.mockController.Finish()
}

type GetSelfHandlerTester struct {
	mockController     *gomock.Controller
	databaseController sqlmock.Sqlmock
	mockedDatabase     *sql.DB
	mockedORM          *bun.DB
	mockedJWT          *mock.MockJWTGenerator
	mockedClock        clock.Clock
	handler            *GetSelfHandler
}

func prepareGetSelfHandlerForTest(t *testing.T) (*GetSelfHandlerTester, error) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockedDatabase, databaseController, err := sqlmock.New()
	if err != nil {
		return nil, err
	}
	databaseController.
		ExpectQuery("SELECT version()").
		WithoutArgs().
		WillReturnRows(sqlmock.NewRows([]string{"version()"}).AddRow("8.0.28"))
	mockedORM := bun.NewDB(mockedDatabase, mysqldialect.New())

	mockedJWT := mock.NewMockJWTGenerator(mockController)

	mockedClock := clock.NewMock()

	dependency := GetSelfHandlerDependencies{
		DB: mockedORM,
	}

	handler, err := NewGetSelfHandler(dependency)
	if err != nil {
		return nil, err
	}

	return &GetSelfHandlerTester{
		mockController:     mockController,
		databaseController: databaseController,
		mockedDatabase:     mockedDatabase,
		mockedORM:          mockedORM,
		mockedJWT:          mockedJWT,
		mockedClock:        mockedClock,
		handler:            handler,
	}, nil
}
