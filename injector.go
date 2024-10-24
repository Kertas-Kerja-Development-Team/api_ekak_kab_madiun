//go:build wireinject
// +build wireinject

package main

import (
	"ekak_kabupaten_madiun/app"
	"ekak_kabupaten_madiun/controller"
	"ekak_kabupaten_madiun/middleware"
	"ekak_kabupaten_madiun/repository"
	"ekak_kabupaten_madiun/service"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"github.com/julienschmidt/httprouter"
)

var rencanaKinerjaSet = wire.NewSet(
	repository.NewRencanaKinerjaRepositoryImpl,
	wire.Bind(new(repository.RencanaKinerjaRepository), new(*repository.RencanaKinerjaRepositoryImpl)),
	service.NewRencanaKinerjaServiceImpl,
	wire.Bind(new(service.RencanaKinerjaService), new(*service.RencanaKinerjaServiceImpl)),
	controller.NewRencanaKinerjaControllerImpl,
	wire.Bind(new(controller.RencanaKinerjaController), new(*controller.RencanaKinerjaControllerImpl)),
)

var rencanaAksiSet = wire.NewSet(
	repository.NewRencanaAksiRepositoryImpl,
	wire.Bind(new(repository.RencanaAksiRepository), new(*repository.RencanaAksiRepositoryImpl)),
	service.NewRencanaAksiServiceImpl,
	wire.Bind(new(service.RencanaAksiService), new(*service.RencanaAksiServiceImpl)),
	controller.NewRencanaAksiControllerImpl,
	wire.Bind(new(controller.RencanaAksiController), new(*controller.RencanaAksiControllerImpl)),
)

var pelaksanaanRencanaAksiSet = wire.NewSet(
	repository.NewPelaksanaanRencanaAksiRepositoryImpl,
	wire.Bind(new(repository.PelaksanaanRencanaAksiRepository), new(*repository.PelaksanaanRencanaAksiRepositoryImpl)),
	service.NewPelaksanaanRencanaAksiServiceImpl,
	wire.Bind(new(service.PelaksanaanRencanaAksiService), new(*service.PelaksanaanRencanaAksiServiceImpl)),
	controller.NewPelaksanaanRencanaAksiControllerImpl,
	wire.Bind(new(controller.PelaksanaanRencanaAksiController), new(*controller.PelaksanaanRencanaAksiControllerImpl)),
)

var usulanMusrebangSet = wire.NewSet(
	repository.NewUsulanMusrebangRepositoryImpl,
	wire.Bind(new(repository.UsulanMusrebangRepository), new(*repository.UsulanMusrebangRepositoryImpl)),
	service.NewUsulanMusrebangServiceImpl,
	wire.Bind(new(service.UsulanMusrebangService), new(*service.UsulanMusrebangServiceImpl)),
	controller.NewUsulanMusrebangControllerImpl,
	wire.Bind(new(controller.UsulanMusrebangController), new(*controller.UsulanMusrebangControllerImpl)),
)

var usulanMandatoriSet = wire.NewSet(
	repository.NewUsulanMandatoriRepositoryImpl,
	wire.Bind(new(repository.UsulanMandatoriRepository), new(*repository.UsulanMandatoriRepositoryImpl)),
	service.NewUsulanMandatoriServiceImpl,
	wire.Bind(new(service.UsulanMandatoriService), new(*service.UsulanMandatoriServiceImpl)),
	controller.NewUsulanMandatoriControllerImpl,
	wire.Bind(new(controller.UsulanMandatoriController), new(*controller.UsulanMandatoriControllerImpl)),
)

var usulanPokokPikiranSet = wire.NewSet(
	repository.NewUsulanPokokPikiranRepositoryImpl,
	wire.Bind(new(repository.UsulanPokokPikiranRepository), new(*repository.UsulanPokokPikiranRepositoryImpl)),
	service.NewUsulanPokokPikiranServiceImpl,
	wire.Bind(new(service.UsulanPokokPikiranService), new(*service.UsulanPokokPikiranServiceImpl)),
	controller.NewUsulanPokokPikiranControllerImpl,
	wire.Bind(new(controller.UsulanPokokPikiranController), new(*controller.UsulanPokokPikiranControllerImpl)),
)

var usulanInisiatifSet = wire.NewSet(
	repository.NewUsulanInisiatifRepositoryImpl,
	wire.Bind(new(repository.UsulanInisiatifRepository), new(*repository.UsulanInisiatifRepositoryImpl)),
	service.NewUsulanInisiatifServiceImpl,
	wire.Bind(new(service.UsulanInisiatifService), new(*service.UsulanInisiatifServiceImpl)),
	controller.NewUsulanInisiatifControllerImpl,
	wire.Bind(new(controller.UsulanInisiatifController), new(*controller.UsulanInisiatifControllerImpl)),
)

var usulanTerpilihSet = wire.NewSet(
	repository.NewUsulanTerpilihRepositoryImpl,
	wire.Bind(new(repository.UsulanTerpilihRepository), new(*repository.UsulanTerpilihRepositoryImpl)),
	service.NewUsulanTerpilihServiceImpl,
	wire.Bind(new(service.UsulanTerpilihService), new(*service.UsulanTerpilihServiceImpl)),
	controller.NewUsulanTerpilihControllerImpl,
	wire.Bind(new(controller.UsulanTerpilihController), new(*controller.UsulanTerpilihControllerImpl)),
)

var gambaranUmumSet = wire.NewSet(
	repository.NewGambaranUmumRepositoryImpl,
	wire.Bind(new(repository.GambaranUmumRepository), new(*repository.GambaranUmumRepositoryImpl)),
	service.NewGambaranUmumServiceImpl,
	wire.Bind(new(service.GambaranUmumService), new(*service.GambaranUmumServiceImpl)),
	controller.NewGambaranUmumControllerImpl,
	wire.Bind(new(controller.GambaranUmumController), new(*controller.GambaranUmumControllerImpl)),
)

var dasarHukumSet = wire.NewSet(
	repository.NewDasarHukumRepositoryImpl,
	wire.Bind(new(repository.DasarHukumRepository), new(*repository.DasarHukumRepositoryImpl)),
	service.NewDasarHukumServiceImpl,
	wire.Bind(new(service.DasarHukumService), new(*service.DasarHukumServiceImpl)),
	controller.NewDasarHukumControllerImpl,
	wire.Bind(new(controller.DasarHukumController), new(*controller.DasarHukumControllerImpl)),
)

var inovasiSet = wire.NewSet(
	repository.NewInovasiRepositoryImpl,
	wire.Bind(new(repository.InovasiRepository), new(*repository.InovasiRepositoryImpl)),
	service.NewInovasiServiceImpl,
	wire.Bind(new(service.InovasiService), new(*service.InovasiServiceImpl)),
	controller.NewInovasiControllerImpl,
	wire.Bind(new(controller.InovasiController), new(*controller.InovasiControllerImpl)),
)

var subKegiatanSet = wire.NewSet(
	repository.NewSubKegiatanRepositoryImpl,
	wire.Bind(new(repository.SubKegiatanRepository), new(*repository.SubKegiatanRepositoryImpl)),
	service.NewSubKegiatanServiceImpl,
	wire.Bind(new(service.SubKegiatanService), new(*service.SubKegiatanServiceImpl)),
	controller.NewSubKegiatanControllerImpl,
	wire.Bind(new(controller.SubKegiatanController), new(*controller.SubKegiatanControllerImpl)),
)

var subKegiatanTerpilihSet = wire.NewSet(
	repository.NewSubKegiatanTerpilihRepositoryImpl,
	wire.Bind(new(repository.SubKegiatanTerpilihRepository), new(*repository.SubKegiatanTerpilihRepositoryImpl)),
	service.NewSubKegiatanTerpilihServiceImpl,
	wire.Bind(new(service.SubKegiatanTerpilihService), new(*service.SubKegiatanTerpilihServiceImpl)),
	controller.NewSubKegiatanTerpilihControllerImpl,
	wire.Bind(new(controller.SubKegiatanTerpilihController), new(*controller.SubKegiatanTerpilihControllerImpl)),
)

func InitializeServer() *http.Server {

	wire.Build(
		app.GetConnection,
		wire.Value([]validator.Option{}),
		validator.New,
		rencanaKinerjaSet,
		rencanaAksiSet,
		pelaksanaanRencanaAksiSet,
		usulanMusrebangSet,
		usulanMandatoriSet,
		usulanPokokPikiranSet,
		usulanInisiatifSet,
		usulanTerpilihSet,
		gambaranUmumSet,
		dasarHukumSet,
		inovasiSet,
		subKegiatanSet,
		subKegiatanTerpilihSet,
		app.NewRouter,
		wire.Bind(new(http.Handler), new(*httprouter.Router)),
		middleware.NewAuthMiddleware,
		NewServer,
	)

	return nil
}
