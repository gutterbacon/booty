package registry

import (
	"go.amplifyedge.org/booty-v2/config"
	"go.amplifyedge.org/booty-v2/dep"
	"go.amplifyedge.org/booty-v2/dep/components"
	"go.amplifyedge.org/booty-v2/internal/store"
	"go.amplifyedge.org/booty-v2/internal/update"
)

type Registry struct {
	DevComponents map[string]dep.Component
	Components    map[string]dep.Component
}

func NewRegistry(db *store.DB, ac *config.AppConfig) (*Registry, error) {
	var err error
	// protoc deps
	protoGenGo := components.NewProtocGenGo(db)
	protoGenGrpc := components.NewProtocGenGoGrpc(db)
	protoCobra := components.NewProtocGenCobra(db)

	comps := []dep.Component{
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
	}

	// register it
	devComponents := map[string]dep.Component{} // dev components
	regComponents := map[string]dep.Component{} // regular components
	for _, c := range comps {
		if c.Dependencies() != nil {
			for _, d := range c.Dependencies() {
				if err = setVersion(ac, d); err != nil {
					return nil, err
				}
			}
		}
		if err = setVersion(ac, c); err != nil {
			return nil, err
		}
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

func setVersion(ac *config.AppConfig, c dep.Component) error {
	var err error
	v := ac.GetVersion(c.Name())
	if v == "" {
		v, err = update.GetLatestVersion(c.RepoUrl())
		if err != nil {
			return err
		}
	}
	c.SetVersion(update.Version(v))
	return nil
}
