package dep

import (
	"go.amplifyedge.org/booty-v2/config"
	"go.amplifyedge.org/booty-v2/dep/components"
	"go.amplifyedge.org/booty-v2/internal/store"
	"go.amplifyedge.org/booty-v2/internal/update"
)

type Registry struct {
	DevComponents map[string]Component
	Components    map[string]Component
}

func NewRegistry(db *store.DB, ac *config.AppConfig) (*Registry, error) {
	var err error
	// protoc deps
	protoGenGo := components.NewProtocGenGo(db)
	protoGenGrpc := components.NewProtocGenGoGrpc(db)
	protoCobra := components.NewProtocGenCobra(db)

	comps := []Component{
		components.NewGoreleaser(db),
		components.NewCaddy(db),
		components.NewGrafana(db),
		components.NewProtoc(
			db,
			[]Component{
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
	devComponents := map[string]Component{} // dev components
	regComponents := map[string]Component{} // regular components
	for _, c := range comps {
		if c.IsDev() {
			v := ac.GetVersion(c.Name())
			if v == "" {
				v, err = update.GetLatestVersion(c.RepoUrl())
				if err != nil {
					return nil, err
				}
			}
			c.SetVersion(update.Version(v))
			devComponents[c.Name()] = c
		} else {
			regComponents[c.Name()] = c
		}
	}
	return &Registry{
		DevComponents: devComponents,
		Components:    regComponents,
	}, err
}
