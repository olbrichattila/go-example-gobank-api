package types

import "flag"

type App struct {
	IsSqlite  bool
	IsSeeding bool
	Port      string
}

func (a *App) Init() {
	isSeeding := flag.Bool("seed", false, "seed the db")
	isSqlite := flag.Bool("sqlite", false, "SqLite server")
	port := flag.String("port", "8000", "port to serve")
	flag.Parse()

	a.IsSeeding = *isSeeding
	a.IsSqlite = *isSqlite
	a.Port = ":" + *port
}
