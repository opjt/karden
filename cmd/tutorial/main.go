package main

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// 1. kubeconfig 로드 (KUBECONFIG 환경변수 → ~/.kube/config 순서로 자동 탐색)
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		loadingRules,
		&clientcmd.ConfigOverrides{},
	).ClientConfig()
	if err != nil {
		panic(err)
	}

	// 2. 클라이언트 생성
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// 3. Informer Factory 생성 (모든 네임스페이스)
	factory := informers.NewSharedInformerFactory(clientset, 0)

	// 4. Pod Informer 생성
	podInformer := factory.Core().V1().Pods().Informer()

	// 5. 이벤트 핸들러 등록
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			fmt.Printf("[ADD] %s/%s\n", pod.Namespace, pod.Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			pod := newObj.(*corev1.Pod)
			fmt.Printf("[UPDATE] %s/%s\n", pod.Namespace, pod.Name)
		},
		DeleteFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			fmt.Printf("[DELETE] %s/%s\n", pod.Namespace, pod.Name)
		},
	})

	// 6. Informer 시작
	stopCh := make(chan struct{})
	defer close(stopCh)

	factory.Start(stopCh)
	factory.WaitForCacheSync(stopCh)

	fmt.Println("Pod watching 시작...")

	// 7. 종료 대기
	<-stopCh
}
