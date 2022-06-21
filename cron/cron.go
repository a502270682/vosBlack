package cron

import (
	"context"
	"fmt"
	"time"
	"vosBlack/adapter/log"
	"vosBlack/adapter/logic"
	"vosBlack/model"

	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

var cronCtl Ctl

type Ctl struct {
	C *cron.Cron
}

func init() {
	cronCtl.C = cron.New(cron.WithSeconds(), cron.WithLocation(time.Local))
}

type EnterpriseApplyHourListJob struct {
	JobName string
}

func (s *EnterpriseApplyHourListJob) Run() {
	// do something
	ctx := context.Background()
	enterprises, err := model.GetEnterpriseInfoImpl().GetAllActiveEnterprise(ctx)
	if err != nil {
		log.Warnf(ctx, "fail to get active enterprises, err:%+v", err)
		return
	}
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)
	success := 0
	var failEid []int
	for _, enterprise := range enterprises {
		field, err := logic.GetApplyHourListCache(ctx, enterprise.NID)
		if err != nil || field == nil {
			continue
		}
		eID := enterprise.NID
		err = model.GetApplyHourListImpl().GetDBForTransaction().Transaction(func(tx *gorm.DB) error {
			entity := &model.EnterpriseApplyHourList{
				EnID:           eID,
				DayReport:      now,
				RepYear:        now.Year(),
				RepMonth:       int(now.Month()),
				RepDay:         now.Day(),
				RepHour:        now.Hour(),
				MbRequestCount: field.MbRequestCount,
				MbHitCount:     field.MbHitCount,
				WnHitCount:     field.WnHitCount,
				MpRequestCount: field.MpRequestCount,
				MpHitCount:     field.MpHitCount,
				GwRequestCount: field.GwRequestCount,
				GwHitCount:     field.GwHitCount,
				FqRequestCount: field.FqRequestCount,
				FqHitCount:     field.FqHitCount,
				JoinDt:         now,
				Remark:         field.Remark,
			}
			err = tx.Table("t_enterprise_applyhourlist").Create(entity).Error
			if err != nil {
				return errors.Wrap(err, "upsert applyhourlist failed")
			}
			err = logic.DeleteApplyHourListCache(ctx, eID)
			if err != nil {
				return errors.Wrap(err, "fail to DeleteApplyHourListCache")
			}
			return nil
		})
		if err != nil {
			failEid = append(failEid, eID)
			log.Warnf(ctx, "fail to create applyhourlist, enterprise_id:%d, err:%+v", eID, err)
			continue
		}
		success++
	}
	log.Infof(ctx, "[applyhourlistJob] time:%+v success:%d failEid:%+v", time.Now(), success, failEid)
}

type GwApplyHourListJob struct {
	JobName string
}

func (s *GwApplyHourListJob) Run() {
	// do something
	ctx := context.Background()
	activeGateWays, err := model.GetSysGatewayImpl().GetAllActiveGateWay(ctx)
	if err != nil {
		log.Warnf(ctx, "fail to get active gateway, err:%+v", err)
		return
	}
	now := time.Now()
	success := 0
	var failId []int
	for _, gw := range activeGateWays {
		field, err := logic.GetGwApplyHourListCache(ctx, gw.NID)
		if err != nil || field == nil {
			continue
		}
		gwID := gw.NID
		err = model.GetGateWayApplyHourListImpl().GetDBForTransaction().Transaction(func(tx *gorm.DB) error {
			entity := &model.GatewayApplyHourList{
				GwID:           gwID,
				DayReport:      now,
				RepYear:        now.Year(),
				RepMonth:       int(now.Month()),
				RepDay:         now.Day(),
				RepHour:        now.Hour(),
				MbRequestCount: field.MbRequestCount,
				MbHitCount:     field.MbHitCount,
				JoinDt:         now,
				Remark:         field.Remark,
			}
			err = tx.Table("sys_gateway_applyhourlist").Create(entity).Error
			if err != nil {
				return errors.Wrap(err, "upsert gwapplyhourlist failed")
			}
			err = logic.DeleteGwApplyHourListCache(ctx, gwID)
			if err != nil {
				return errors.Wrap(err, "fail to DeleteGwApplyHourListCache")
			}
			return nil
		})
		if err != nil {
			failId = append(failId, gwID)
			log.Warnf(ctx, "fail to create gwapplyhourlist, gw_id:%d, err:%+v", gwID, err)
			continue
		}
		success++
	}
	log.Infof(ctx, "[gwapplyhourlistJob] time:%+v success:%d failId:%+v", time.Now(), success, failId)
}

func StartCron() {
	if cronCtl.C != nil {
		ctl := cronCtl.C
		enJob := EnterpriseApplyHourListJob{
			JobName: "enterprise_apply_hour_list_job",
		}
		// every 1 hour
		_, err := ctl.AddJob("0 0 0/1 * * *", &enJob)
		if err != nil {
			panic(errors.Wrap(err, fmt.Sprintf("failed to start %s", enJob.JobName)))
		}
		ctl.Start()
	} else {
		panic("cron not init")
	}
	if cronCtl.C != nil {
		ctl := cronCtl.C
		gwJob := GwApplyHourListJob{
			JobName: "gw_apply_hour_list_Job",
		}
		// every 1 hour
		_, err := ctl.AddJob("0 0 0/1 * * *", &gwJob)
		if err != nil {
			panic(errors.Wrap(err, fmt.Sprintf("failed to start %s", gwJob.JobName)))
		}
		ctl.Start()
	} else {
		panic("cron not init")
	}
}

func StopCron() {
	if cronCtl.C != nil {
		cronCtl.C.Stop()
	}
}
