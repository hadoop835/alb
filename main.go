package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"time"

	"k8s.io/klog"

	"alauda.io/alb2/config"
	"alauda.io/alb2/controller"
	"alauda.io/alb2/driver"
	"alauda.io/alb2/ingress"
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()
	defer klog.Flush()
	klog.Error("Service start.")

	err := config.ValidateConfig()
	if err != nil {
		klog.Error(err.Error())
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config.Set("LABEL_SERVICE_ID", fmt.Sprintf("service.%s/uuid", config.Get("DOMAIN")))
	config.Set("LABEL_SERVICE_NAME", fmt.Sprintf("service.%s/name", config.Get("DOMAIN")))
	config.Set("LABEL_CREATOR", fmt.Sprintf("service.%s/createby", config.Get("DOMAIN")))

	d, err := driver.GetDriver()
	if err != nil {
		panic(err)
	}
	// install necessary crd on start
	if config.GetBool("INSTALL_CRD") {
		if err := d.RegisterCustomDefinedResources(); err != nil {
			// install crd failed, abort
			panic(err)
		}
	}

	go ingress.MainLoop(ctx)
	go func() {
		// for profiling
		http.ListenAndServe(":1937", nil)
	}()

	if config.Get("LB_TYPE") == config.Nginx {
		go rotateLog()
	}

	interval := config.GetInt("INTERVAL")
	tmo := time.Duration(config.GetInt("RELOAD_TIMEOUT")) * time.Second
	for {
		time.Sleep(time.Duration(interval) * time.Second)
		ch := make(chan string)
		startTime := time.Now()
		klog.Info("Begin update reload loop")

		go func() {
			err := controller.TryLockAlb()
			if err != nil {
				klog.Error("lock alb failed", err.Error())
			}
			ctl, err := controller.GetController()
			if err != nil {
				klog.Error(err.Error())
				ch <- "continue"
				return
			}
			ch <- "wait"

			ctl.GC()
			err = ctl.GenerateConf()
			if err != nil {
				klog.Error(err.Error())
				ch <- "continue"
				return
			}
			ch <- "wait"
			err = ctl.ReloadLoadBalancer()
			if err != nil {
				klog.Error(err.Error())
			}
			ch <- "continue"
			return
		}()
		timer := time.NewTimer(tmo)

	watchdog:
		for {
			select {
			case msg := <-ch:
				if msg == "continue" {
					klog.Info("continue")
					timer.Reset(0)
					break watchdog
				}
				timer.Reset(tmo)
				continue
			case <-timer.C:
				klog.Error("reload timeout")
				klog.Flush()
				os.Exit(1)
			}
		}

		klog.Infof("End update reload loop, cost %s", time.Since(startTime))
	}
}

func rotateLog() {
	rotateInterval := config.GetInt("ROTATE_INTERVAL")
	klog.Info("rotateLog start, rotate interval ", rotateInterval)
	for {
		klog.Info("start rorate log")
		output, err := exec.Command("/usr/sbin/logrotate", "/etc/logrotate.d/alauda").CombinedOutput()
		if err != nil {
			klog.Errorf("rotate log failed %s %v", output, err)
		} else {
			klog.Info("rotate log success")
		}
		time.Sleep(time.Duration(rotateInterval) * time.Minute)
	}
}
