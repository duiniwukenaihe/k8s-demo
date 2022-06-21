package service

import (
	coreV1 "k8s.io/api/core/v1"
	"strconv"
	"strings"
)

func (d *Deployment) GetLabels() map[string]string {
	labelsMap := make(map[string]string)
	labels := strings.Split(d.Labels, "\n")
	for _, label := range labels {
		values := strings.SplitN(label, ":", 2)
		if len(values) != 2 {
			continue
		}
		labelsMap[strings.TrimSpace(values[0])] = strings.TrimSpace(values[1])
	}
	return labelsMap
}

func (d *Deployment) GetSelectors() map[string]string {
	selectors := d.GetLabels()
	selectors["app"] = d.Name
	return selectors
}
func (d *Deployment) GetPorts() []coreV1.ContainerPort {
	portList := make([]coreV1.ContainerPort, 0, 5)
	ports := strings.Split(d.Ports, "\n")
	for _, port := range ports {
		values := strings.SplitN(port, ",", 3)
		if len(values) != 3 {
			continue
		}
		intPort, err := strconv.Atoi(values[1])
		if err != nil {
			continue
		}
		protocol := coreV1.ProtocolTCP
		if strings.Compare(strings.ToLower(values[0]), "tcp") != 0 {
			protocol = coreV1.ProtocolUDP
		}
		portList = append(portList, coreV1.ContainerPort{
			Name:          strings.TrimSpace(values[2]),
			ContainerPort: int32(intPort),
			Protocol:      protocol,
		})
	}

	return portList
}
func (d *Deployment) GetImageName() string {
	// 全部为应为字母数字和:
	pods := strings.Index(d.Images, ":")
	if pods > 0 {
		return d.Images[:pods]
	}
	return d.Images
}
