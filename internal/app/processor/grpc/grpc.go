package pgrpc

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/KDarenskii/catalog-service/internal/app/config/section"
	"github.com/KDarenskii/catalog-service/internal/app/processor"
	catalogv1 "github.com/KDarenskii/catalog-service/internal/pkg/grpc/gen/catalog/v1"
)

type grpcProc struct {
	server *grpc.Server
	addr   string
}

func NewGRPC(
	catalogV1 catalogv1.CatalogServiceServer,
	cfg section.ProcessorGrpc,
) processor.Processor {
	srv := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))

	catalogv1.RegisterCatalogServiceServer(srv, catalogV1)

	reflection.Register(srv)

	return &grpcProc{server: srv, addr: fmt.Sprintf(":%d", cfg.ListenPort)}
}

func (p *grpcProc) StartAsync(ctx context.Context, wg *sync.WaitGroup) {
	var lc net.ListenConfig

	l, err := lc.Listen(ctx, "tcp", p.addr)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start GRPC server")
		return
	}

	log.Info().Str("address", p.addr).Str("network", "TCP").Msg("Started GRPC server")

	go func() {
		if err := p.server.Serve(l); err != nil {
			log.Fatal().Err(err).Msg("Failed to start GRPC server")
			return
		}
	}()

	processor.WatchForShutdown(
		ctx, wg, processor.CloserFunc(
			func() error {
				p.server.GracefulStop()
				return nil
			},
		),
	)
}
