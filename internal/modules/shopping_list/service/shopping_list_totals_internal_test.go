package service

import (
	"testing"

	"github.com/stretchr/testify/require"

	shoppingModel "github.com/nclsgg/despensa-digital/backend/internal/modules/shopping_list/model"
)

func TestCalculateListTotalsHandlesMassConversion(t *testing.T) {
	items := []shoppingModel.ShoppingListItem{
		{
			Quantity:       300, // grams
			Unit:           "g",
			EstimatedPrice: 30, // price per kilogram
			Purchased:      true,
			ActualPrice:    28,
		},
	}

	estimated, actual := calculateListTotals(items)
	require.InEpsilon(t, 9, estimated, 1e-6)
	require.InEpsilon(t, 8.4, actual, 1e-6)
}

func TestCalculateListTotalsHandlesVolumeConversion(t *testing.T) {
	items := []shoppingModel.ShoppingListItem{
		{
			Quantity:       750,
			Unit:           "ml",
			EstimatedPrice: 20, // per liter
		},
		{
			Quantity:       2,
			Unit:           "un",
			EstimatedPrice: 5,
			Purchased:      true,
			ActualPrice:    4.5,
		},
	}

	estimated, actual := calculateListTotals(items)
	require.InEpsilon(t, 20*0.75+10, estimated, 1e-6)
	require.InEpsilon(t, 4.5*2, actual, 1e-6)
}
