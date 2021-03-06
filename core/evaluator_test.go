package core

import (
	"github.com/stretchr/testify/assert"
	"github.com/PerfLibNetcracker/gorse/base"
	"math"
	"strconv"
	"testing"
)

const evalEpsilon = 0.00001

func EqualEpsilon(t *testing.T, expect float64, actual float64, epsilon float64) {
	if math.Abs(expect-actual) > evalEpsilon {
		t.Fatalf("Expect %f ± %f, Actual: %f\n", expect, epsilon, actual)
	}
}

type EvaluatorTesterModel struct {
	MaxUserId int
	MaxItemId int
	Matrix    [][]float64
}

func NewEvaluatorTesterModel(users, items []int, ratings []float64) *EvaluatorTesterModel {
	test := new(EvaluatorTesterModel)
	if len(users) == 0 {
		test.MaxUserId = -1
	} else {
		test.MaxUserId = base.Max(users)
	}
	if len(items) == 0 {
		test.MaxItemId = -1
	} else {
		test.MaxItemId = base.Max(items)
	}
	test.Matrix = base.NewMatrix(test.MaxUserId+1, test.MaxItemId+1)
	for i := range ratings {
		userId := users[i]
		itemId := items[i]
		rating := ratings[i]
		test.Matrix[userId][itemId] = rating
	}
	return test
}

func (tester *EvaluatorTesterModel) Predict(userId, itemId string) float64 {
	userIndex, err1 := strconv.Atoi(userId)
	itemIndex, err2 := strconv.Atoi(itemId)
	if err1 != nil || err2 != nil {
		return 0
	} else if userIndex > tester.MaxUserId {
		return 0
	} else if itemIndex > tester.MaxItemId {
		return 0
	} else {
		return tester.Matrix[userIndex][itemIndex]
	}
}

func (tester *EvaluatorTesterModel) GetParams() base.Params {
	panic("EvaluatorTesterModel.GetParams() should never be called.")
}

func (tester *EvaluatorTesterModel) SetParams(params base.Params) {
	panic("EvaluatorTesterModel.SetParams() should never be called.")
}

func (tester *EvaluatorTesterModel) Fit(set DataSetInterface, options *base.RuntimeOptions) {
	panic("EvaluatorTesterModel.Fit() should never be called.")
}

func NewTestIndexer(id []string) *base.Indexer {
	indexer := base.NewIndexer()
	for _, v := range id {
		indexer.Add(v)
	}
	return indexer
}

func TestRMSE(t *testing.T) {
	a := []float64{0, 0, 0}
	b := []float64{-2.0, 0, 2.0}
	if math.Abs(RMSE(a, b)-1.63299) > evalEpsilon {
		t.Fail()
	}
}

func TestMAE(t *testing.T) {
	a := []float64{0, 0, 0}
	b := []float64{-2.0, 0, 2.0}
	if math.Abs(MAE(a, b)-1.33333) > evalEpsilon {
		t.Fail()
	}
}

func NewTestTargetSet(ids []string) *base.MarginalSubSet {
	indexer := NewTestIndexer(ids)
	indices := make([]int, len(ids))
	values := make([]float64, len(ids))
	subset := make([]int, len(ids))
	for i := range subset {
		subset[i] = i
		indices[i] = i
	}
	return base.NewMarginalSubSet(indexer, indices, values, subset)
}

func TestNDCG(t *testing.T) {
	targetSet := NewTestTargetSet([]string{"1", "3", "5", "7"})
	rankList := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	EqualEpsilon(t, 0.6766372989, NDCG(targetSet, rankList), evalEpsilon)
}

func TestPrecision(t *testing.T) {
	targetSet := NewTestTargetSet([]string{"1", "3", "5", "7"})
	rankList := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	EqualEpsilon(t, 0.4, Precision(targetSet, rankList), evalEpsilon)
}

func TestRecall(t *testing.T) {
	targetSet := NewTestTargetSet([]string{"1", "3", "15", "17", "19"})
	rankList := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	EqualEpsilon(t, 0.4, Recall(targetSet, rankList), evalEpsilon)
}

func TestAP(t *testing.T) {
	targetSet := NewTestTargetSet([]string{"1", "3", "7", "9"})
	rankList := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	EqualEpsilon(t, 0.44375, MAP(targetSet, rankList), evalEpsilon)
}

func TestRR(t *testing.T) {
	targetSet := NewTestTargetSet([]string{"3"})
	rankList := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	EqualEpsilon(t, 0.25, MRR(targetSet, rankList), evalEpsilon)
}

func TestAUC(t *testing.T) {
	// The mocked test dataset:
	// 1.0 0.0 0.0
	// 0.0 0.5 0.0
	// 0.0 0.0 1.0
	a := NewEvaluatorTesterModel([]int{0, 0, 0, 1, 1, 1, 2, 2, 2},
		[]int{0, 1, 2, 0, 1, 2, 0, 1, 2},
		[]float64{1.0, 0.0, 0.0, 0.0, 0.5, 0.0, 0.0, 0.0, 1.0})
	b := NewDataSet([]string{"0", "1", "2"}, []string{"0", "1", "2"}, []float64{1.0, 0.5, 1.0})
	assert.Equal(t, 1.0, EvaluateAUC(a, b, nil))
}
