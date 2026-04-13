package ambulance_wl

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TechOctopus/davgus-ambulance-webapi/internal/db_service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type PatientsApiSuite struct {
	suite.Suite
	patientDbMock   *DbServiceMock[Patient]
	placementDbMock *DbServiceMock[Placement]
}

type DbServiceMock[DocType any] struct {
	mock.Mock
}

func (m *DbServiceMock[DocType]) CreateDocument(ctx context.Context, id string, document *DocType) error {
	args := m.Called(ctx, id, document)
	return args.Error(0)
}

func (m *DbServiceMock[DocType]) FindDocument(ctx context.Context, id string) (*DocType, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*DocType), args.Error(1)
}

func (m *DbServiceMock[DocType]) FindAllDocuments(ctx context.Context) ([]DocType, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]DocType), args.Error(1)
}

func (m *DbServiceMock[DocType]) UpdateDocument(ctx context.Context, id string, document *DocType) error {
	args := m.Called(ctx, id, document)
	return args.Error(0)
}

func (m *DbServiceMock[DocType]) DeleteDocument(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *DbServiceMock[DocType]) Disconnect(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

var _ db_service.DbService[Patient] = (*DbServiceMock[Patient])(nil)
var _ db_service.DbService[Placement] = (*DbServiceMock[Placement])(nil)

func TestPatientsApiSuite(t *testing.T) {
	suite.Run(t, new(PatientsApiSuite))
}

func (suite *PatientsApiSuite) SetupTest() {
	suite.patientDbMock = &DbServiceMock[Patient]{}
	suite.placementDbMock = &DbServiceMock[Placement]{}
}

func (suite *PatientsApiSuite) Test_DeletePatient_ArchivesAndRemovesPlacements() {
	// ARRANGE
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service_patients", suite.patientDbMock)
	ctx.Set("db_service_placements", suite.placementDbMock)
	ctx.Params = []gin.Param{{Key: "patientId", Value: "test-patient"}}

	suite.patientDbMock.
		On("FindDocument", mock.Anything, "test-patient").
		Return(&Patient{Id: "test-patient", Name: "Test", Archived: false}, nil)
	suite.patientDbMock.
		On("UpdateDocument", mock.Anything, "test-patient", mock.Anything).
		Return(nil)

	suite.placementDbMock.
		On("FindAllDocuments", mock.Anything).
		Return([]Placement{
			{Id: "pl-1", PatientId: "test-patient"},
			{Id: "pl-2", PatientId: "other-patient"},
		}, nil)
	suite.placementDbMock.
		On("DeleteDocument", mock.Anything, "pl-1").
		Return(nil)

	sut := implPatientsAPI{}

	// ACT
	sut.DeletePatient(ctx)

	// ASSERT
	suite.Equal(http.StatusNoContent, recorder.Code)
	suite.patientDbMock.AssertCalled(suite.T(), "FindDocument", mock.Anything, "test-patient")
	suite.patientDbMock.AssertCalled(
		suite.T(),
		"UpdateDocument",
		mock.Anything,
		"test-patient",
		mock.MatchedBy(func(document *Patient) bool {
			return document != nil && document.Archived
		}),
	)
	suite.placementDbMock.AssertCalled(suite.T(), "DeleteDocument", mock.Anything, "pl-1")
	suite.placementDbMock.AssertNotCalled(suite.T(), "DeleteDocument", mock.Anything, "pl-2")
}

func (suite *PatientsApiSuite) Test_GetPatients_ReturnsOnlyActive() {
	// ARRANGE
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service_patients", suite.patientDbMock)

	suite.patientDbMock.
		On("FindAllDocuments", mock.Anything).
		Return([]Patient{
			{Id: "p-active", Name: "Active", Archived: false},
			{Id: "p-archived", Name: "Archived", Archived: true},
		}, nil)

	sut := implPatientsAPI{}

	// ACT
	sut.GetPatients(ctx)

	// ASSERT
	suite.Equal(http.StatusOK, recorder.Code)

	var body []Patient
	err := json.Unmarshal(recorder.Body.Bytes(), &body)
	suite.Require().NoError(err)
	suite.Len(body, 1)
	suite.Equal("p-active", body[0].Id)
	suite.False(body[0].Archived)
}
