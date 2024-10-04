package ticketoption

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"ticket-api/mocks"
	"ticket-api/storage"
	"time"
)

type TicketOptionTestSuite struct {
	suite.Suite
	StorageMock *mocks.MockTicketOptionStorage
	Controller  *gomock.Controller
}

func (s *TicketOptionTestSuite) SetupSuite() {
	s.Controller = gomock.NewController(s.T())
	s.StorageMock = mocks.NewMockTicketOptionStorage(s.Controller)
}

func (s *TicketOptionTestSuite) TestHandleAllocationNotAvailable() {
	storageInput := &storage.GenerateTicketsInput{
		UserId:         "chris-test",
		TicketOptionId: "1234",
		Quantity:       15,
	}

	s.StorageMock.EXPECT().GenerateTickets(gomock.Any(), storageInput).Return(nil, storage.ErrTicketOptionAllocationCheckFailed)

	defaultTicketOptionService := NewDefaultTicketOptionService(s.StorageMock)
	err := defaultTicketOptionService.PurchaseTicketOption(context.TODO(), &PurchaseTicketOptionInput{
		Quantity:       15,
		UserID:         "chris-test",
		TicketOptionId: "1234",
	})

	assert.ErrorIs(s.T(), err, ErrOverAllocatedTickets)
}

func (s *TicketOptionTestSuite) TestCanPurchaseTickets() {
	storageInput := &storage.GenerateTicketsInput{
		UserId:         "chris-test",
		TicketOptionId: "1234",
		Quantity:       3,
	}

	s.StorageMock.EXPECT().GenerateTickets(gomock.Any(), storageInput).Return(&storage.GenerateTicketsResult{
		PurchaseId: "1234",
		TicketIds:  []string{"5", "6", "7"},
	}, nil)

	defaultTicketOptionService := NewDefaultTicketOptionService(s.StorageMock)
	err := defaultTicketOptionService.PurchaseTicketOption(context.TODO(), &PurchaseTicketOptionInput{
		Quantity:       3,
		UserID:         "chris-test",
		TicketOptionId: "1234",
	})

	assert.Nil(s.T(), err)
}

func (s *TicketOptionTestSuite) TestNotEnoughTicketsGenerated() {
	storageInput := &storage.GenerateTicketsInput{
		UserId:         "chris-test",
		TicketOptionId: "1234",
		Quantity:       5,
	}

	s.StorageMock.EXPECT().GenerateTickets(gomock.Any(), storageInput).Return(&storage.GenerateTicketsResult{
		PurchaseId: "1234",
		TicketIds:  []string{"5", "6", "7"},
	}, nil)

	defaultTicketOptionService := NewDefaultTicketOptionService(s.StorageMock)
	err := defaultTicketOptionService.PurchaseTicketOption(context.TODO(), &PurchaseTicketOptionInput{
		Quantity:       5,
		UserID:         "chris-test",
		TicketOptionId: "1234",
	})

	assert.ErrorIs(s.T(), err, ErrNotEnoughTicketsGenerated)
}

func (s *TicketOptionTestSuite) TestCreateTicketOptionSuccess() {
	s.StorageMock.EXPECT().CreateTicketOption(gomock.Any(), &storage.CreateTicketOptionInput{
		Name:        "my test event",
		Description: "a test event",
		Allocation:  100,
	}).Return(&storage.TicketOptionResult{
		ID:          "1234",
		Name:        "my test event",
		Description: "a test event",
		Allocation:  100,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	}, nil)

	defaultTicketOptionService := NewDefaultTicketOptionService(s.StorageMock)
	to, err := defaultTicketOptionService.CreateTicketOption(context.TODO(), &CreateTicketOptionInput{
		Name:        "my test event",
		Description: "a test event",
		Allocation:  100,
	})

	expectedTo := &TicketOption{
		ID:          "1234",
		Name:        "my test event",
		Description: "a test event",
		Allocation:  100,
	}

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expectedTo, to)
}

func TestTicketOptionTestSuite(t *testing.T) {
	suite.Run(t, new(TicketOptionTestSuite))
}
