package pgateway

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/KDarenskii/catalog-service/internal/app/config/section"
	"github.com/KDarenskii/catalog-service/internal/app/processor"
	catalogv1 "github.com/KDarenskii/catalog-service/internal/pkg/grpc/gen/catalog/v1"
)

const (
	gatewayReadHeaderTimeout = 5 * time.Second
	gatewayReadTimeout       = 30 * time.Second
	gatewayWriteTimeout      = 30 * time.Second
	gatewayIdleTimeout       = 120 * time.Second
	gatewayShutdownTimeout   = 5 * time.Second
)

type gatewayProc struct {
	server   http.Server
	addr     string
	grpcAddr string
}

func NewGateway(
	cfgGateway section.ProcessorGateway,
	cfgGrpc section.ProcessorGrpc,
) processor.Processor {
	address := fmt.Sprintf(":%d", cfgGateway.ListenPort)
	grpcAddress := net.JoinHostPort("localhost", strconv.Itoa(int(cfgGrpc.ListenPort)))

	return &gatewayProc{
		server: http.Server{
			ReadTimeout:       gatewayReadTimeout,
			ReadHeaderTimeout: gatewayReadHeaderTimeout,
			WriteTimeout:      gatewayWriteTimeout,
			IdleTimeout:       gatewayIdleTimeout,
		},
		addr:     address,
		grpcAddr: grpcAddress,
	}
}

func (p *gatewayProc) StartAsync(ctx context.Context, wg *sync.WaitGroup) {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := catalogv1.RegisterCatalogServiceHandlerFromEndpoint(ctx, mux, p.grpcAddr, opts); err != nil {
		log.Fatal().Err(err).Msg("Failed to register gateway catalog handler")
		return
	}

	p.server.Handler = mux

	var lc net.ListenConfig

	l, err := lc.Listen(ctx, "tcp", p.addr)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start gateway server")
	}

	log.Info().Str("address", p.addr).Str("network", "TCP").Msg("Started gateway server")

	go func() {
		if err := p.server.Serve(l); err != nil {
			log.Fatal().Err(err).Msg("Failed to start gateway server")
			return
		}
	}()

	processor.WatchForShutdown(ctx, wg, processor.CloserFunc(l.Close))

	processor.WatchForShutdown(ctx, wg, processor.NewCloserContextFunc(p.server.Shutdown, context.Background(), 5*time.Second))
}
