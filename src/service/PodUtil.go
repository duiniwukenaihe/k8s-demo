package service

import (
	corev1 "k8s.io/api/core/v1"

	"strconv"
	"strings"
)

func GetPodMessage(pod corev1.Pod) string {
	message := ""
	for _, contition := range pod.Status.Conditions {
		if contition.Status != "True" {
			message += contition.Message
		}
	}
	return message
}
func GetPodIsReady(pod corev1.Pod) bool {
	for _, condition := range pod.Status.Conditions {
		if condition.Type == "ContainersReady" && condition.Status != "True" {
			return false
		}
	}

	for _, rg := range pod.Spec.ReadinessGates {
		for _, condition := range pod.Status.Conditions {
			if condition.Type == rg.ConditionType && condition.Status != "True" {
				return false
			}
		}
	}
	return true
}
func (d *Pod) GetPorts() []corev1.ContainerPort {
	portList := make([]corev1.ContainerPort, 0, 5)
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
		protocol := corev1.ProtocolTCP
		if strings.Compare(strings.ToLower(values[0]), "tcp") != 0 {
			protocol = corev1.ProtocolUDP
		}
		portList = append(portList, corev1.ContainerPort{
			Name:          strings.TrimSpace(values[2]),
			ContainerPort: int32(intPort),
			Protocol:      protocol,
		})
	}

	return portList
}
func (d *Pod) GetImageName() string {
	// 全部为应为字母数字和:
	pods := strings.Index(d.Images, ":")
	if pods > 0 {
		return d.Images[:pods]
	}
	return d.Images
}
