package menu

import (
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type ProjectMenu struct {
	Id         int64
	Pid        int64
	Title      string
	Icon       string
	Url        string
	FilePath   string
	Role       string
	Status     int
	Values     string
	ShowSlider int
}

func (*ProjectMenu) TableName() string {
	return "menu"
}

type ProjectMenuChild struct {
	ProjectMenu
	Children []*ProjectMenuChild
}

func CovertChild(pms []*ProjectMenu) []*ProjectMenuChild {
	var pmcs []*ProjectMenuChild
	err := copier.Copy(&pmcs, pms)
	if err != nil {
		zap.L().Error("ProjectMenu模块结构体赋值错误", zap.Error(err))
		return nil
	}
	var childPmcs []*ProjectMenuChild
	//递归
	for _, v := range pmcs {
		if v.Pid == 0 {
			pmc := &ProjectMenuChild{}
			err := copier.Copy(pmc, v)
			if err != nil {
				zap.L().Error("ProjectMenu模块结构体嵌套赋值错误", zap.Error(err))
			}
			childPmcs = append(childPmcs, pmc)
		}
	}
	toChild(childPmcs, pmcs)
	return childPmcs
}

func toChild(childPmcs []*ProjectMenuChild, pmcs []*ProjectMenuChild) {
	for _, pmc := range childPmcs {
		for _, pm := range pmcs {
			if pmc.Id == pm.Pid {
				child := &ProjectMenuChild{}
				err := copier.Copy(child, pm)
				if err != nil {
					zap.L().Error("ProjectMenu结构体赋值错误", zap.Error(err))
				}
				pmc.Children = append(pmc.Children, child)
			}
		}
		toChild(pmc.Children, pmcs)
	}
}
