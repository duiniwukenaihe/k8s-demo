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
