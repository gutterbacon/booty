package registry

import (
	"go.amplifyedge.org/booty-v2/dep"
	"go.amplifyedge.org/booty-v2/dep/components"
	"go.amplifyedge.org/booty-v2/internal/store"
)

type Registry struct {
	DevComponents map[string]dep.Component
	Components    map[string]dep.Component
}

func NewRegistry(db store.Storer) (*Registry, error) {
	var err error
	// protoc deps
	protoGenGo := components.NewProtocGenGo(db)
	protoGenGrpc := components.NewProtocGenGoGrpc(db)
	protoCobra := components.NewProtocGenCobra(db)

	comps := []dep.Component{
		components.NewBooty(db),
		components.NewGoreleaser(db),
		components.NewCaddy(db),
		components.NewGrafana(db),
		components.NewProtoc(
			db,
			[]dep.Component{
				protoGenGo,
				protoGenGrpc,
				protoCobra,
			},
		),
		components.NewGoJsonnet(db),
		components.NewVicMet(db),
		components.NewJb(db),
		components.NewMkcert(db),
		components.NewProtocGenInjectTag(db),
	}

	// register it
	devComponents := map[string]dep.Component{} // dev components
	regComponents := map[string]dep.Component{} // regular components
	for _, c := range comps {
		if !c.IsDev() {
			regComponents[c.Name()] = c
		}
		devComponents[c.Name()] = c
	}
	return &Registry{
		DevComponents: devComponents,
		Components:    regComponents,
	}, err
}
