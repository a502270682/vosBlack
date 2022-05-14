package model

type IStatus int

const (
	IStatusActive IStatus = 1
	IStatusStop   IStatus = 0
	IStatusPause  IStatus = 9
	IStatusDelete IStatus = -1
)

type GwType int

const (
	GwTypeVosBlack        GwType = 0
	GwTypeVosRewrite      GwType = 1
	GwTypeYunxuntongBlack GwType = 2
	GwTypeDongyunBlack    GwType = 3
	GwTypeHuaxinVosBlack  GwType = 4
)
