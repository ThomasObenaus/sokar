package main

import gomock "github.com/golang/mock/gomock"

type dummyReceiver struct {
}

func (dr dummyReceiver) Any() {
}

func newGomockCall(ctrl *gomock.Controller) *gomock.Call {
	return ctrl.RecordCall(dummyReceiver{}, "Any")
}
