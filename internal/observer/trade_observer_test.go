package observer

import (
	"testing"

	"TickFlow/internal/database"

	"go.uber.org/mock/gomock"
)

func TestTradeObserver_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := database.NewMockDatabase(ctrl)

	testData := map[string]interface{}{
		"symbol": "AAPL",
		"price":  150.0,
	}

	mockDB.EXPECT().SaveTrade(testData).Return(nil)

	observer := NewTradeObserver(mockDB)

	observer.Update(testData)
}
