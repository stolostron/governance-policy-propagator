package propagator

import (
	"fmt"
	"sync"
	"time"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var logger = logf.Log.WithName("izhang batchUpdater")

type batchUpdater struct {
	// mutex for the set
	mtx *sync.Mutex
	set map[reconcile.Request]struct{}
	// reconciler called by the manager might conflict with us
	// we redo it 2mins later
	rec reconcile.Reconciler

	defaultUpdateInterval time.Duration
}

func NewbatchUpdater(r reconcile.Reconciler) *batchUpdater {
	return &batchUpdater{
		mtx: &sync.Mutex{},
		rec: r,
		set: map[reconcile.Request]struct{}{},

		defaultUpdateInterval: time.Second * 7,
	}
}

// every to 2mins clean up the set by adding all the items to controller queue
func (b *batchUpdater) Start(stop <-chan struct{}) error {
	go b.run(stop)

	logger.Info("Started")
	return nil
}

func (b *batchUpdater) run(stop <-chan struct{}) error {
	slowTicker := time.NewTicker(b.defaultUpdateInterval)
	defer slowTicker.Stop()

	for {
		select {
		case <-stop:
			logger.Info("Stopped")
			return nil
		case <-slowTicker.C:
			b.update()

		}
	}
}

func (b *batchUpdater) add(req reconcile.Request) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	if _, ok := b.set[req]; !ok {
		b.set[req] = struct{}{}
	}
}

func (b *batchUpdater) update() {
	logger.Info(fmt.Sprintf("enter update() %v policies ", len(b.set)))
	defer logger.Info("exit update()")

	pro_set := []reconcile.Request{}

	b.mtx.Lock()

	for k := range b.set {
		pro_set = append(pro_set, k)
	}

	b.mtx.Unlock()

	for _, req := range pro_set {
		logger.Info(fmt.Sprintf("policy: %s ", req))
		if _, err := b.rec.Reconcile(req); err != nil {
			continue
		}

		delete(b.set, req)
	}

	b.mtx.Lock()
	for _, req := range pro_set {
		b.set[req] = struct{}{}
	}

	b.mtx.Unlock()
}
