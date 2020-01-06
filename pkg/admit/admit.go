package admit

import (
	"time"

	"github.com/oam-dev/oam-go-sdk/pkg/client/informers/externalversions/core.oam.dev/v1alpha1"

	"github.com/oam-dev/oam-go-sdk/pkg/client/clientset/versioned"
	"github.com/oam-dev/oam-go-sdk/pkg/client/informers/externalversions"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// Admit is the main object for admission controller
type Admit struct {
	Client            *versioned.Clientset
	Factory           externalversions.SharedInformerFactory
	componentInformer v1alpha1.ComponentSchematicInformer
	traitInformer     v1alpha1.TraitInformer
	scopeInformer     v1alpha1.ApplicationScopeInformer
	appConfigInformer v1alpha1.ApplicationConfigurationInformer
}

// New is the entrance of getting an Admit object
func New() (*Admit, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	client, err := versioned.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	factory := externalversions.NewSharedInformerFactory(client, 6*time.Hour)
	componentInformer := factory.Core().V1alpha1().ComponentSchematics()
	componentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{})
	traitInformer := factory.Core().V1alpha1().Traits()
	traitInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{})
	scopeInformer := factory.Core().V1alpha1().ApplicationScopes()
	scopeInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{})
	appConfigInformer := factory.Core().V1alpha1().ApplicationConfigurations()
	appConfigInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{})

	return &Admit{Client: client,
		Factory:           factory,
		componentInformer: componentInformer,
		traitInformer:     traitInformer,
		scopeInformer:     scopeInformer,
		appConfigInformer: appConfigInformer,
	}, nil
}

func (adm *Admit) Start(stop <-chan struct{}) {
	adm.Factory.Start(stop)
	cache.WaitForCacheSync(stop, adm.appConfigInformer.Informer().HasSynced)
	cache.WaitForCacheSync(stop, adm.componentInformer.Informer().HasSynced)
	cache.WaitForCacheSync(stop, adm.traitInformer.Informer().HasSynced)
	cache.WaitForCacheSync(stop, adm.scopeInformer.Informer().HasSynced)
}
