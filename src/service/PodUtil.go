package service

import v1 "k8s.io/api/core/v1"

func GetPodMessage(pod v1.Pod) string {
	message := ""
	for _, contition := range pod.Status.Conditions {
		if contition.Status != "True" {
			message += contition.Message
		}
	}
	return message
}
func GetPodIsReady(pod v1.Pod) bool {
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
