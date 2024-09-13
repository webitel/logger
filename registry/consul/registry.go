package consul

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/webitel/logger/model"
	"github.com/webitel/logger/registry"
	"log/slog"
	"net"
	"strconv"
	"time"
)

type ConsulRegistry struct {
	registrationConfig *consulapi.AgentServiceRegistration
	client             *consulapi.Client
	stop               chan any
	checkId            string
}

func NewConsulRegistry(config *model.ConsulConfig) (*ConsulRegistry, model.AppError) {
	var (
		err error
	)
	entity := ConsulRegistry{}
	if config.Id == "" {
		return nil, model.NewBadRequestError("consul.registry.new_consul.check_args.service_id", "service id is empty! (set it by '-id' flag)")
	}
	ip, port, err := net.SplitHostPort(config.PublicAddress)
	if err != nil {
		return nil, model.NewBadRequestError("consul.registry.new_consul.parse_address.error", "unable to parse address")
	}
	parsedPort, err := strconv.Atoi(port)
	if err != nil {
		return nil, model.NewBadRequestError("consul.registry.new_consul.parse_ip.error", "unable to parse ip")
	}

	consulConfig := consulapi.DefaultConfig()
	consulConfig.Address = config.Address
	entity.client, err = consulapi.NewClient(consulConfig)
	if err != nil {
		return nil, model.NewBadRequestError("consul.registry.new_consul_registry.consulapi_creation.error", err.Error())
	}

	entity.registrationConfig = &consulapi.AgentServiceRegistration{
		ID:      config.Id,
		Name:    registry.ServiceName,
		Port:    parsedPort,
		Address: ip,
		Check: &consulapi.AgentServiceCheck{
			DeregisterCriticalServiceAfter: registry.DeregisterCriticalServiceAfter.String(),
			CheckID:                        config.Id,
			TTL:                            registry.CheckInterval.String(),
		},
	}
	entity.stop = make(chan any)

	return &entity, nil
}

func (c *ConsulRegistry) Register() model.AppError {
	agent := c.client.Agent()
	err := agent.ServiceRegister(c.registrationConfig)
	if err != nil {
		return model.NewInternalError("consul.registry.consul.call_register.error", err.Error())
	}
	var checks map[string]*consulapi.AgentCheck
	if checks, err = agent.Checks(); err != nil {
		return model.NewInternalError("consul.registry.consul.register.get_checks.error", err.Error())
	}

	var serviceCheck *consulapi.AgentCheck
	for _, check := range checks {
		if check.ServiceID == c.registrationConfig.ID {
			serviceCheck = check
		}
	}

	if serviceCheck == nil {
		return model.NewInternalError("consul.registry.consul.register.error", err.Error())
	}
	c.checkId = serviceCheck.CheckID
	go c.RunServiceCheck()
	slog.Info(fmtConsulLog("service was registered"))
	return nil
}

func (c *ConsulRegistry) Deregister() model.AppError {
	err := c.client.Agent().ServiceDeregister(c.registrationConfig.ID)
	if err != nil {
		return model.NewInternalError("consul.registry.consul.register.error", err.Error())
	}
	c.stop <- true
	slog.Info(fmtConsulLog("service was deregistered"))
	return nil
}

func (c *ConsulRegistry) RunServiceCheck() model.AppError {
	defer slog.Info(fmtConsulLog("stopped service checker"))
	slog.Info(fmtConsulLog("started service checker"))
	ticker := time.NewTicker(registry.CheckInterval / 2)
	for {
		select {
		case <-c.stop:
			// gracefull stop
			return nil
		case <-ticker.C:
			err := c.client.Agent().UpdateTTL(c.checkId, "success", "pass")
			if err != nil {
				slog.Warn(fmtConsulLog(err.Error()))
			}
			// TODO: seems that connection is lost, reconnect?
		}
	}
}

func fmtConsulLog(s string) string {
	return fmt.Sprintf("consul: %s", s)
}
