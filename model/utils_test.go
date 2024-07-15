package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github/shivasaicharanruthala/backend-engineer-assignment/errors"
)

func TestIsValidUUID(t *testing.T) {
	testCase := []struct {
		id           int
		useCase      string
		uuidStr      string
		expectedResp bool
	}{
		{
			id: 1, useCase: "Negative case: invalid uuid",
			uuidStr:      "1278",
			expectedResp: false,
		},
		{
			id: 2, useCase: "Negative case: empty uuid",
			uuidStr:      "",
			expectedResp: false,
		},
		{
			id: 3, useCase: "Positive case: valid uuid",
			uuidStr:      "4a77ec9d-5334-43d0-a9e1-4fca8807bf8f",
			expectedResp: true,
		},
	}

	for _, tc := range testCase {
		resp := IsValidUUID(tc.uuidStr)
		assert.Equal(t, tc.expectedResp, resp, fmt.Sprintf("Test %v Failed with use case %v", tc.id, tc.useCase))
	}
}

func TestStringPointer(t *testing.T) {
	expectedStr1 := ""
	expectedStr2 := "fetch rewards"

	testCase := []struct {
		id           int
		useCase      string
		inpStr       string
		expectedResp *string
	}{
		{
			id: 1, useCase: "Positive case: empty string",
			inpStr:       "",
			expectedResp: &expectedStr1,
		},
		{
			id: 2, useCase: "Positive case: string with finite length",
			inpStr:       "fetch rewards",
			expectedResp: &expectedStr2,
		},
	}

	for _, tc := range testCase {
		resp := StringPointer(tc.inpStr)
		assert.Equal(t, tc.expectedResp, resp, fmt.Sprintf("Test %v Failed with use case %v", tc.id, tc.useCase))
	}
}

func TestSixPointRule(t *testing.T) {
	testCase := []struct {
		id          int
		useCase     string
		receipt     *Receipt
		expectedErr error
	}{
		{
			id: 1, useCase: "Negative case: invalid purchaseDate",
			receipt:     &Receipt{PurchaseDate: StringPointer("05-02")},
			expectedErr: errors.NewInvalidParam(errors.InvalidParam{Param: "purchaseDate"}),
		},
		{
			id: 2, useCase: "Positive case: valid purchaseDate with odd day",
			receipt:     &Receipt{PurchaseDate: StringPointer("2022-05-01")},
			expectedErr: nil,
		},
		{
			id: 3, useCase: "Positive case: valid purchaseDate with even day",
			receipt:     &Receipt{PurchaseDate: StringPointer("2022-05-02")},
			expectedErr: nil,
		},
	}

	for _, tc := range testCase {
		err := tc.receipt.SixPointRule()
		if err != nil {
			assert.Equal(t, tc.expectedErr.Error(), err.Error(), fmt.Sprintf("Test %v Failed with use case %v", tc.id, tc.useCase))
		} else {
			assert.Equal(t, tc.expectedErr, err, fmt.Sprintf("Test %v Failed with use case %v", tc.id, tc.useCase))
		}
	}
}

func TestTenPointRule(t *testing.T) {
	testCase := []struct {
		id          int
		useCase     string
		receipt     *Receipt
		expectedErr error
	}{
		{
			id: 1, useCase: "Negative case: invalid purchaseTime",
			receipt:     &Receipt{PurchaseTime: StringPointer("02")},
			expectedErr: errors.NewInvalidParam(errors.InvalidParam{Param: "purchaseTime"}),
		},
		{
			id: 2, useCase: "Negative case: invalid purchaseTime, hours greater than 24",
			receipt:     &Receipt{PurchaseTime: StringPointer("25:02")},
			expectedErr: errors.NewInvalidParam(errors.InvalidParam{Param: "purchaseTime"}),
		},
		{
			id: 3, useCase: "Negative case: invalid purchaseTime, minutes greater than 60",
			receipt:     &Receipt{PurchaseTime: StringPointer("15:61")},
			expectedErr: errors.NewInvalidParam(errors.InvalidParam{Param: "purchaseTime"}),
		},
		{
			id: 4, useCase: "Positive case: valid purchaseTime, before 2pm",
			receipt:     &Receipt{PurchaseTime: StringPointer("13:00")},
			expectedErr: nil,
		},
		{
			id: 5, useCase: "Positive case: valid purchaseTime, after 4pm",
			receipt:     &Receipt{PurchaseTime: StringPointer("17:00")},
			expectedErr: nil,
		},
		{
			id: 6, useCase: "Positive case: valid purchaseTime, between 2pm to 4pm",
			receipt:     &Receipt{PurchaseTime: StringPointer("15:00")},
			expectedErr: nil,
		},
	}

	for _, tc := range testCase {
		err := tc.receipt.TenPointRule()
		if err != nil {
			assert.Equal(t, tc.expectedErr.Error(), err.Error(), fmt.Sprintf("Test %v Failed with use case %v", tc.id, tc.useCase))
		} else {
			assert.Equal(t, tc.expectedErr, err, fmt.Sprintf("Test %v Failed with use case %v", tc.id, tc.useCase))
		}
	}
}

func TestCalculateTotalReceiptPoints(t *testing.T) {
	testCase := []struct {
		id             int
		useCase        string
		receipt        *Receipt
		expectedPoints int
		expectedErr    error
	}{
		{
			id: 1, useCase: "New Receipt: All rules are covered",
			receipt: &Receipt{
				Retailer:     StringPointer("M&M Corner Market 45"), // Retailer has Letters, Numbers, Special Chars
				PurchaseDate: StringPointer("2022-01-01"),
				PurchaseTime: StringPointer("15:01"),
				Total:        StringPointer("25.00"),
				Items: []Item{
					{ShortDescription: StringPointer("Mountain Dew 12PK"), Price: StringPointer("5.00")},
					{ShortDescription: StringPointer("Emils Cheese Pizza"), Price: StringPointer("10.00")},
					{ShortDescription: StringPointer("Doritos Nacho Cheese"), Price: StringPointer("5.00")},
					{ShortDescription: StringPointer("   Klarbrunn 12-PK 12 FL OZ  "), Price: StringPointer("5.00")},
				},
			},
			expectedPoints: 120,
		},
		{
			id: 2, useCase: "New Receipt: Retailer has only letters, day is even, PurchaseTime is not between 2 to 4pm, Total has cents",
			receipt: &Receipt{
				Retailer:     StringPointer("Target"),
				PurchaseDate: StringPointer("2022-01-02"),
				PurchaseTime: StringPointer("13:01"),
				Total:        StringPointer("5.05"),
				Items: []Item{
					{ShortDescription: StringPointer("Mountain Dew 12PK"), Price: StringPointer("5.05")},
				},
			},
			expectedPoints: 6,
		},
		{
			id: 3, useCase: "Negative case: Invalid year",
			receipt: &Receipt{
				Retailer:     StringPointer("Target"),
				PurchaseDate: StringPointer("20-01-02"),
				PurchaseTime: StringPointer("13:01"),
				Total:        StringPointer("5.05"),
				Items: []Item{
					{ShortDescription: StringPointer("Mountain Dew 12PK"), Price: StringPointer("5.05")},
				},
			},
			expectedErr: errors.NewInvalidParam(errors.InvalidParam{Param: "purchaseDate"}),
		},
		{
			id: 3, useCase: "Negative case: Invalid time",
			receipt: &Receipt{
				Retailer:     StringPointer("Target"),
				PurchaseDate: StringPointer("2024-01-02"),
				PurchaseTime: StringPointer(":01"),
				Total:        StringPointer("5.05"),
				Items: []Item{
					{ShortDescription: StringPointer("Mountain Dew 12PK"), Price: StringPointer("5.05")},
				},
			},
			expectedErr: errors.NewInvalidParam(errors.InvalidParam{Param: "purchaseTime"}),
		},
	}

	for _, tc := range testCase {
		err := tc.receipt.CalculateTotalReceiptPoints()
		if err != nil {
			assert.Equal(t, tc.expectedErr.Error(), err.Error(), fmt.Sprintf("Test %v Failed with use case %v", tc.id, tc.useCase))
		} else {
			assert.Equal(t, tc.expectedPoints, tc.receipt.Points, fmt.Sprintf("Test %v Failed with use case %v", tc.id, tc.useCase))
			assert.Equal(t, tc.expectedErr, err, fmt.Sprintf("Test %v Failed with use case %v", tc.id, tc.useCase))
		}
	}
}
